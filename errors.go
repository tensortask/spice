package spice

import "errors"

var (
	// ErrAlreadyClosed ...
	ErrAlreadyClosed = errors.New("spice: encryptor already closed")

	// ErrInvalidData ...
	ErrInvalidData = errors.New("spice: encrypted message is invalid")
)
