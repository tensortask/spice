package spice

import (
	"io"

	"golang.org/x/crypto/nacl/box"
)

// Decryptor is an io.ReadCloser that reads encrypted data written by an
// Encryptor.
type Decryptor struct {
	reader     io.Reader
	nonce      *[24]byte
	sharedKey  *[32]byte
	chunkSize  int
	inputSlice []byte
}

// NewDecryptor returns a new Decryptor. Nonce and key should be identical to
// the values originally passed to NewEncryptor.
//
// Neither nonce or key are modified.
func NewDecryptor(reader io.Reader, nonce *[24]byte, sharedKey *[32]byte) *Decryptor {
	return &Decryptor{reader: reader, nonce: nonce, sharedKey: sharedKey, chunkSize: DefaultChunkSize}
}

// Decrypt ///
func (d *Decryptor) Decrypt(outputSlice []byte) error {
	bytesRead, err := d.reader.Read(d.inputSlice)
	if err != nil {
		return err
	}
	outputSlice, _ = box.OpenAfterPrecomputation(outputSlice, d.inputSlice[:bytesRead], d.nonce, d.sharedKey)
	copy(d.nonce[:], outputSlice[:24])
	return nil
}
