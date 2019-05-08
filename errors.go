package spice

import "errors"

var (
	// ErrOpenNotOk ...
	ErrOpenNotOk = errors.New("spice: box.OpenAfterPrecomputation returned a NOT ok bool (most-likely verification failure)")
)
