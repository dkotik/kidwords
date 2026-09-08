/*
Package internal contains common project helpers and utilities.
*/
package internal

import (
	"io"
	"math/rand/v2"
)

func NewDetermenisticRandomReader(seed [32]byte) io.Reader {
	return rand.NewChaCha8(seed)
}
