package spice

import "io"

// Decryptor is an io.ReadCloser that reads encrypted data written by an
// Encryptor.
type Decryptor struct {
	r     io.Reader
	nonce *[24]byte
	key   *[32]byte
}

// NewDecryptor returns a new Decryptor. Nonce and key should be identical to
// the values originally passed to NewEncryptor.
//
// Neither nonce or key are modified.
func NewDecryptor(r io.Reader, nonce *[16]byte, key *[32]byte) (*Decryptor, error) {
	d := Decryptor{r: r}

	return &d, nil
}
