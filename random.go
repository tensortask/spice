package spice

import (
	"crypto/rand"
	"io"
)

// RandomNonce ...
func RandomNonce() ([24]byte, error) {
	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return [24]byte{}, err
	}
	return nonce, nil
}
