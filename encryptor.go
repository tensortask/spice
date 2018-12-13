package spice

import (
	"io"

	"golang.org/x/crypto/nacl/box"
)

// Encryptor is an io.WriteCloser. Writes to an Encryptor are encrypted
// and written to w.
type Encryptor struct {
	writer      io.Writer
	nonce       *[24]byte
	sharedKey   *[32]byte
	chunkSize   int
	outputSlice []byte
}

// NewEncryptor creates an Encryptor with the default chunk size.
func NewEncryptor(writer io.Writer, nonce *[24]byte, sharedKey *[32]byte) *Encryptor {
	return &Encryptor{writer: writer, nonce: nonce, sharedKey: sharedKey, chunkSize: DefaultChunkSize}
}

// Encrypt /..
func (e *Encryptor) Encrypt(inputSlice []byte) error {
	e.outputSlice = box.SealAfterPrecomputation(nil, inputSlice, e.nonce, e.sharedKey)
	_, err := e.writer.Write(e.outputSlice)
	if err != nil {
		return err
	}
	copy(e.nonce[:], e.outputSlice[:24])
	return nil
}
