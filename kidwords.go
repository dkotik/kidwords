/*
Package kidwords provides data encoding accessible to children.

Use it to create passwords and paper keys.

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
