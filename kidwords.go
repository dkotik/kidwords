/*
Package kidwords provides durable and accessible paper key encoding that children can use.

Printable paper keys are occasionally used as the last resort for recovering account access. They increase security by empowering a user with the ability to wrestle control of a compromised account from an attacker.

Most paper keys are encoded using BIP39 convention into a set of words. The final few words encode the integrity of the key with a cyclical redundancy check. When printed and stored, such keys are not durable because they can be lost to minor physical damage.

Kid Words package or command line tool increases key durability by splitting the key using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) algorithm into shards and encoding each shard using a dictionary of 256 four-letter English nouns.

## Benefits

- Keys can be recovered from partially damaged paper.
- Shards can be transmitted and memorized by children.
- Shards are easier to speak over poor radio or telephone connection, which can save time during an emergency.
- Key shards can be hidden in several physical locations by cutting the paper into pieces. Once a configurable quorum of shards, three by default, is gathered back, the key can be restored.
- Shards can easily be obfuscated by sequencing:
  - toys or books on a shelf
  - pencil scribbles on paper
  - objects or signs in a Minecraft world
  - emojis

- Command line tool can apply all of the above benefits to:
  - important passwords to rarely accessed accounts that do not support paper keys
  - conventional BIP39 keys

## Inspired By

- [Horcrux][horcrux]

horcrux: https://github.com/jesseduffield/horcrux/tree/master
*/
package kidwords

import (
	"bytes"
	"io"
	"strings"
)

//go:generate go run dictionary/generate.go --source dictionary/enNouns.txt --destination dictionary/enNouns.gen.go --variable EnglishFourLetterNouns
//go:generate go test . -update

// FromReader translates [io.Reader] stream into Kid Words.
func FromReader(r io.Reader, withOptions ...WriterOption) (string, error) {
	b := bytes.Buffer{}
	w, err := NewWriter(&b, withOptions...)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(w, r); err != nil {
		return "", err
	}
	return b.String(), nil
}

// FromBytes translates a set of bytes into Kid Words.
func FromBytes(b []byte, withOptions ...WriterOption) (string, error) {
	return FromReader(bytes.NewReader(b), withOptions...)
}

// FromString translates a string into Kid Words.
func FromString(s string, withOptions ...WriterOption) (string, error) {
	return FromReader(strings.NewReader(s), withOptions...)
}

// func FromInt(r io.Reader) (n int64, err error) {
// 	wrapped := NewReader(r)
// 	b := &bytes.Buffer{}
// 	_, err = io.CopyN(b, wrapped, 24)
// 	if err != nil && err != io.EOF {
// 		return 0, err
// 	}
//
// 	chopped, ok := ChecksumChop(b.Bytes())
// 	if !ok {
// 		return 0, errors.New("checksum did not match")
// 	}
//
// 	x := new(big.Int).SetBytes(chopped)
// 	return x.Int64(), nil
// }

// ToWriter streams translated Kid Words into [io.Writer].
func ToWriter(w io.Writer, s string, withOptions ...ReaderOption) error {
	r, err := NewReader(strings.NewReader(s), withOptions...)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

// ToBytes translates Kid Words into bytes.
func ToBytes(s string, withOptions ...ReaderOption) ([]byte, error) {
	b := bytes.Buffer{}
	if err := ToWriter(&b, s, withOptions...); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// ToString translates Kid Words into a string.
func ToString(s string, withOptions ...ReaderOption) (string, error) {
	b := bytes.Buffer{}
	if err := ToWriter(&b, s, withOptions...); err != nil {
		return "", err
	}
	return b.String(), nil
}
