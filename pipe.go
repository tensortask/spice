package spice

import (
	"io"

	"golang.org/x/crypto/nacl/box"
)

// EncyptedPipe ...
type EncyptedPipe struct {
	r          io.Reader
	w          io.Writer
	publicKey  *[32]byte
	privateKey *[32]byte
}

// NewEncyptedPipe ...
func NewEncyptedPipe(r io.Reader, w io.Writer, publicKey *[32]byte, privateKey *[32]byte) *EncyptedPipe {
	enc := EncyptedPipe{r: r, w: w, publicKey: publicKey, privateKey: privateKey}
	return &enc
}

// Encrypt ...
func (p *EncyptedPipe) Encrypt(nonce *[24]byte, publicKey *[32]byte) error {
	non := *nonce
	// Precomputing the shared key speeds up encryption.
	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, publicKey, p.privateKey)

	for {
		buffer := make([]byte, ChunkSize)
		_, err := p.r.Read(buffer)

		if err == io.EOF {
			break
		}
		encrypted := box.SealAfterPrecomputation(nil, buffer, &non, sharedEncryptKey)
		_, err = p.w.Write(encrypted)
		if err != nil {
			return err
		}
		copy(non[:], encrypted[:24])
	}
	return nil
}

// Decrypt ...
func (p *EncyptedPipe) Decrypt(nonce *[24]byte, publicKey *[32]byte) error {
	non := *nonce
	// Precomputing the shared key speeds up encryption.
	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, publicKey, p.privateKey)
	for {
		buffer := make([]byte, ChunkSize+Overhead)
		_, err := p.r.Read(buffer)
		if err == io.EOF {
			break
		}
		decrypted, _ := box.OpenAfterPrecomputation(nil, buffer, &non, sharedEncryptKey)
		_, err = p.w.Write(decrypted)
		if err != nil {
			return err
		}

		copy(non[:], buffer[:24])
	}
	return nil
}
