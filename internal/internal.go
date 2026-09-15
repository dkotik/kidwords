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

//go:generate go run generate.go --source wordlist/en/nouns.txt --destination ../dictionary_nouns.gen.go --variable EnglishFourLetterNouns
//go:generate go run generate.go --source wordlist/en/verbs.txt --destination ../dictionary_verbs.gen.go --variable EnglishFourLetterVerbs

func NewDetermenisticRandomReader(seed [32]byte) io.Reader {
	return rand.NewChaCha8(seed)
}

func NewDeterministicUUID(seed string) string {
	return uuid.UUID(
		md5.Sum([]byte(seed)),
	).String()
}
