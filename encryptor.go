package spice

import (
	"io"
)

// Encryptor is an io.WriteCloser. Writes to an Encryptor are encrypted
// and written to w.
type Encryptor struct {
	w     io.Writer // underlying writer
	nonce *[24]byte // nacl nonce, increments per chunk
	key   *[32]byte // encryption key
}

// NewEncryptor creates an Encryptor with the default chunk size.
func NewEncryptor(w io.Writer, nonce *[24]byte, key *[32]byte) *Encryptor {
	enc := Encryptor{w: w, nonce: nonce, key: key}
	return &enc
}
