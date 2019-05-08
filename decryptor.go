package spice

import (
	"io"

	"golang.org/x/crypto/nacl/box"
)

// Decryptor is an io.ReadCloser that reads encrypted data written by an
// Encryptor.
type Decryptor struct {
	writer      io.Writer
	nonce       *[24]byte
	sharedKey   *[32]byte
	chunkSize   int
	outputSlice []byte
}

// NewDecryptor returns a new Decryptor. Nonce and key should be identical to
// the values originally passed to NewEncryptor.
//
// Neither nonce or key are modified.
func NewDecryptor(writer io.Writer, nonce [24]byte, sharedKey *[32]byte) *Decryptor {
	return &Decryptor{writer: writer, nonce: &nonce, sharedKey: sharedKey, chunkSize: DefaultChunkSize}
}

// Decrypt ///
func (d *Decryptor) Decrypt(inputSlice []byte) error {
	var ok bool
	d.outputSlice, ok = box.OpenAfterPrecomputation(nil, inputSlice, d.nonce, d.sharedKey)
	if !ok {
		return ErrOpenNotOk
	}
	_, err := d.writer.Write(d.outputSlice)
	if err != nil {
		return err
	}
	copy(d.nonce[:], inputSlice[:24])
	return nil
}
