package pairing

import (
	"github.com/dedis/paper_17_dfinity/pbc"
	"gopkg.in/dedis/crypto.v0/abstract"
)


// XXX non-NIST ciphers?

// SHA256 hash function
/*
func (s *suiteEd25519) Hash() hash.Hash {
	return sha256.New()
}

// SHA3/SHAKE128 Sponge Cipher
func (s *suiteEd25519) Cipher(key []byte, options ...interface{}) abstract.Cipher {
	return sha3.NewShakeCipher128(key, options...)
}

func (s *suiteEd25519) Read(r io.Reader, objs ...interface{}) error {
	return abstract.SuiteRead(s, r, objs)
}

func (s *suiteEd25519) Write(w io.Writer, objs ...interface{}) error {
	return abstract.SuiteWrite(s, w, objs)
}

func (s *suiteEd25519) New(t reflect.Type) interface{} {
	return abstract.SuiteNew(s, t)
}

// NewKey returns a formatted Ed25519 key (avoiding subgroup attack by requiring
// it to be a multiple of 8)
func (s *suiteEd25519) NewKey(stream cipher.Stream) abstract.Scalar {
	if stream == nil {
		stream = random.Stream
	}
	buffer := random.NonZeroBytes(32, stream)
	scalar := sha512.Sum512(buffer)
	scalar[0] &= 0xf8
	scalar[31] &= 0x3f
	scalar[31] |= 0x40

	secret := s.Scalar().SetBytes(scalar[:32])
	return secret
}

*/

var Suite = NewAES128SHA256Ed25519(false)
var Pairing = NewAES128SHA256Ed25519P(false)

// Ciphersuite based on AES-128, SHA-256, and the Ed25519 curve.
func NewAES128SHA256Ed25519(fullGroup bool) (abstract.Suite)  {
	suite := pbc.NewPairingFp254BNb()
	return suite.G2()
}

func NewAES128SHA256Ed25519P(fullGroup bool) *pbc.Pairing {
	suite := pbc.NewPairingFp254BNb()
	return suite
}

