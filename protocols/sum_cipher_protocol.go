package protocols

import (
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
	"errors"
	"gopkg.in/dedis/onet.v1/log"
	"math/big"
	"time"

	"unlynx/prio_utils"

	"github.com/henrycg/prio/utils"

)


const SumCipherProtocolName = "SumCipher"


/*Messages
____________________________________________________________________________________________________________________
 */

//structure to announce start of protocol
type AnnounceSumCipher struct {
}

//Reply from the children
type ReplySumCipherBytes struct {
	Bytes []byte
}

type ReplySumCipherLength struct {
	BigIntLen int
	BitLen int
}

type CorShare struct {
	CorShareD []byte
	CorShareE []byte
}

type OutShare struct {
	Out		[]byte
}
/*Structs
_________________________________________________________________________________________________________________________
*/

type StructAnnounce struct {
	*onet.TreeNode
	AnnounceSumCipher
}


type StructReply struct {
	*onet.TreeNode
	ReplySumCipherBytes
}

type StructCorShare struct {
	*onet.TreeNode
	CorShare
}

type StructOutShare struct {
	*onet.TreeNode
	OutShare
}

type Cipher struct {
	Share *big.Int

	//for the moment put bit in int
	Bits []uint
}


type AcceptReply struct {
}

type SumCipherProtocol struct {
	*onet.TreeNodeInstance

	//the feedback final
	Feedback chan *big.Int

	//Channel for up and down communication
	ChildDataChannel chan []StructReply

	AnnounceChannel chan StructAnnounce

	//The data of the protocol
	Ciphers []Cipher
	Sum 	*big.Int
	Modulus *big.Int

	//for proofs
	Proofs  bool
	Request []*prio_utils.Request
	pre		*prio_utils.CheckerPrecomp
	Checker	*prio_utils.Checker

	//channel for proof
	CorShareChannel	chan []StructCorShare
	OutShareChannel		chan []StructOutShare

}



type StatusFlag int

// Status of a client submission.
const (
	NotStarted    StatusFlag = iota
	OpenedTriples StatusFlag = iota
	Layer1        StatusFlag = iota
	Finished      StatusFlag = iota
)

type RequestStatus struct {
	check *prio_utils.Checker
	flag  StatusFlag
}
/*
_______________________________________________________________________________
 */
var randomKey = utils.RandomPRGKey()

func init() {
	network.RegisterMessage(AnnounceSumCipher{})
	network.RegisterMessage(ReplySumCipherBytes{})
	network.RegisterMessage(CorShare{})
	onet.GlobalProtocolRegister(SumCipherProtocolName,NewSumCipherProtocol)
}


func NewSumCipherProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance,error) {

	st := &SumCipherProtocol{
		TreeNodeInstance: n,
		Feedback:         make(chan *big.Int),
		Sum:              big.NewInt(int64(0)),
	}

	//register the channel for announce
	err := st.RegisterChannel(&st.AnnounceChannel)
	if err != nil {
		return nil, errors.New("couldn't register Announce data channel: " + err.Error())
	}

	//register the channel for child response
	err = st.RegisterChannel(&st.ChildDataChannel)
	if err != nil {
		return nil, errors.New("couldn't register Child Response channel" + err.Error())
	}

	err = st.RegisterChannel(&st.CorShareChannel)
	if err != nil {
		return nil, errors.New("Couldn't register CorrShare channel" + err.Error())
	}

	err = st.RegisterChannel(&st.OutShareChannel)
	if err !=nil {
		return nil,errors.New("Couldn't register OutShare channel" + err.Error())
	}

	return st,nil
}

//start called at the root
func (p*SumCipherProtocol) Start() error {
	if p.Ciphers == nil {
		return errors.New("No Shares to collect")
	}
	log.Lvl1(p.ServerIdentity(), " started a Sum Cipher Protocol (", len(p.Ciphers), " different shares)")

	start := time.Now()
	//send to the children of the root
	p.SendToChildren(&AnnounceSumCipher{})
	log.Lvl1("time to send mesage to children of root ", time.Since(start))
	return nil
}
//dispatch is called on the node and handle incoming messages

func (p*SumCipherProtocol) Dispatch() error {

	//Go down the tree
	if !p.IsRoot() {
		p.sumCipherAnnouncementPhase()
	}

	//Ascending aggreg
	start := time.Now()
	sum := p.ascendingAggregationPhase()
	log.Lvl1(p.ServerIdentity(), " completed aggregation phase (", sum, " is the sum ) in ", time.Since(start))

	//report result
	if p.IsRoot() {
		p.Feedback <-sum
	}
	return nil
}

func (p *SumCipherProtocol) sumCipherAnnouncementPhase() {
	//send down the tree if you have some
	AnnounceMessage := <-p.AnnounceChannel
	if !p.IsLeaf() {
		p.SendToChildren(&AnnounceMessage.AnnounceSumCipher)
	}
}

// Results pushing up the tree containing aggregation results.
func (p *SumCipherProtocol) ascendingAggregationPhase() *big.Int {


	if p.Ciphers == nil {
		p.Sum = big.NewInt(0)
	}

	if !p.IsLeaf() {
		//wait on the channel for child to complete and add sum
		for _, v := range <-p.ChildDataChannel {
			//get the bytes and turn them back in big.Int
			var sum big.Int
			sum.SetBytes(v.Bytes)
			p.Sum.Add(p.Sum, &sum)
			p.Sum.Mod(p.Sum,p.Modulus)
		}
	}

	//do the sum of ciphers

	for _, v := range p.Ciphers {
		if !Verify(v) {
			log.Lvl1("Share refused, will not use it for the operation ")
		} else {
			p.Sum.Add(p.Sum, Decode(v))
			p.Sum.Mod(p.Sum, p.Modulus)
		}
	}

	//send to parent the sum to deblock channel wait
	if !p.IsRoot() {
		//send the big.Int in bytes
		p.SendToParent(&ReplySumCipherBytes{p.Sum.Bytes()})
	}

	//finish by returning the sum of the root
	p.Sum.Mod(p.Sum,p.Modulus)


	if (p.Proofs) {
		status := new(RequestStatus)

		status.check = p.Checker


		log.Lvl1("request before is " , status.check.Prg)
		status.check.SetReq(p.Request[p.Index()])
		log.Lvl1("N in checker is " ,status.check.N)
		log.Lvl1("request is ", status.check)
		log.Lvl1(p.Tree().Size())

		evalReplies := make([]*prio_utils.CorShare, 1)
		//need to do this for all shares so for all servers


		//here evalReplies filled
		evalReplies[0] = status.check.CorShare(p.pre)

		//From here need to wait all evalReplies
		if !p.IsRoot() {
			log.Lvl1("corshare is ",evalReplies[0])
			p.SendTo(p.Root(),&CorShare{evalReplies[0].ShareD.Bytes(),evalReplies[0].ShareE.Bytes()})
		}
		if(p.IsRoot()) {
			p.SendToChildren(&CorShare{evalReplies[0].ShareD.Bytes(),evalReplies[0].ShareE.Bytes()})
		}


		//actually they need to all send shares to each other so can all reconstruct core

		evalRepliesFromAll := make([]*prio_utils.CorShare,1)
		evalRepliesFromAll[0] = evalReplies[0]

		//when 1 share do not work else, wait on nothing
		if(p.Tree().Size()>1) {
			for _, v := range <-p.CorShareChannel {
				corshare := new(prio_utils.CorShare)
				corshare.ShareD = big.NewInt(0).SetBytes(v.CorShareD)
				corshare.ShareE = big.NewInt(0).SetBytes(v.CorShareE)
				evalRepliesFromAll = append(evalRepliesFromAll, corshare)
			}
		}
		log.Lvl1("will fuse corShare on :",evalRepliesFromAll[0], "and",evalRepliesFromAll[1])

		//cor is same for all server you cannot transfer it that's why you transfer the shares
		cor := status.check.Cor(evalRepliesFromAll)

		//we need to do this on all servers
		finalReplies := make([]*prio_utils.OutShare, 1)
		log.Lvl1(randomKey)
		finalReplies[0] = status.check.OutShare(cor, randomKey)
		log.Lvl1("finalReplies should not be the same", finalReplies[0])

		if(!p.IsRoot()) {
			p.SendTo(p.Root(),&OutShare{finalReplies[0].Check.Bytes()})
		}

		if(p.IsRoot()) {
			finalRepliesAll := make([]*prio_utils.OutShare,1)
			finalRepliesAll[0] = finalReplies[0]
			for _, v := range <-p.OutShareChannel {
				outShare := new(prio_utils.OutShare)
				outShare.Check = big.NewInt(0).SetBytes(v.OutShare.Out)
				finalRepliesAll = append(finalRepliesAll,outShare)
			}
			log.Lvl1("will evaluate on ", finalRepliesAll[0].Check , finalRepliesAll[1].Check)
			log.Lvl1("output is valid ? ", status.check.OutputIsValid(finalRepliesAll))
		}

	}
	return p.Sum

}




func Encode(x *big.Int) (Cipher) {
	length := x.BitLen()
	resultBit := make([]uint,length)
	for i := 0; i < length; i++ {
		resultBit[i] = x.Bit(i)
	}
	cipher := Cipher{x,resultBit}
	return cipher
}

func Verify(c Cipher) (bool) {
	verify := big.NewInt(0)
	for i,b := range c.Bits {
		if b>1 || b<0 {
			panic("Not bits form in the encoding")
			return false
		}
		verify.Add(verify,big.NewInt(0).Mul(big.NewInt(int64(b)),big.NewInt(0).Exp(big.NewInt(2),big.NewInt(int64(i)),nil)))

	}
	difference := big.NewInt(int64(0))
	difference.Sub(c.Share,verify)
	if difference.Uint64()== uint64(0) {
		return true
	}
	errors.New(" The share is not equal to it's bit form")
	return false
}

func Decode(c Cipher)(x *big.Int) {
	return c.Share
}

