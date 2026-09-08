/*
Package internal contains common project helpers and utilities.
*/
package internal

import (
	"crypto/md5"
	"io"
	"math/rand/v2"
	"uuid"
)

func NewDetermenisticRandomReader(seed [32]byte) io.Reader {
	return rand.NewChaCha8(seed)
}

func NewDeterministicUUID(seed string) string {
	return uuid.UUID(
		md5.Sum([]byte(seed)),
	).String()
}
