package prio_utils

import (
	"github.com/henrycg/prio/share"
	"github.com/henrycg/prio/triple"

)

type Uuid [32]byte

const DEFAULT_MAX_PENDING_REQS = 64

// The data struct that the client gives to each server.
type ClientRequest struct {
	Hint *share.PRGHints

	// Compressed representation of Beaver triples for the
	// batch checking and for the main MPC protocol.
	TripleShare *triple.Share
}

type ServerCiphertext struct {
	Nonce      [24]byte // NaCl Box nonce
	Ciphertext []byte   // Encrypted upload payload
}

type UploadArgs struct {
	PublicKey   [32]byte // NaCl Box public key
	Ciphertexts []ServerCiphertext
}


type NewRequestArgs struct {
	RequestID  Uuid
	Ciphertext ServerCiphertext
}

type NewRequestReply struct {
}





