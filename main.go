package spice

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/nacl/box"
)

// DefaultChunkSize is the size in bytes of a given encrypted chunk.
//
const DefaultChunkSize int = 16368

// DefaultOverhead ...
const DefaultOverhead int = 16

// Keypair ...
type Keypair struct {
	Public  *[32]byte
	Private *[32]byte
}

// SharedKey ...
func (k *Keypair) SharedKey(peerPublicKey *[32]byte) *[32]byte {
	sharedKey := new([32]byte)
	box.Precompute(sharedKey, peerPublicKey, k.Private)
	return sharedKey
}

// RandomNonce ...
func RandomNonce() (*[24]byte, error) {
	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}
	return &nonce, nil
}

// RandomKeypair ...
func RandomKeypair() (Keypair, error) {
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return Keypair{}, err
	}
	return Keypair{Public: publicKey, Private: privateKey}, nil
}

// NewKeypair ...
func NewKeypair(public *[32]byte, private *[32]byte) Keypair {
	return Keypair{Public: public, Private: private}
	//
}
