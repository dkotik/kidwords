/*
Package dictionary defines arrays of 256 words used for KidWords encoding.
*/
package dictionary

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

// Dictionary holds 256 words, each corresponding to a byte value.
type Dictionary [256]string

func (d *Dictionary) Reverse() map[string]byte {
	m := make(map[string]byte)
	for i, w := range d {
		m[w] = byte(i)
	}
	return m
}

// Validate iterates through every value to check for uniqueness and extra white space characters.
func (d *Dictionary) Validate() error {
	if d == nil {
		return errors.New("provided dictionary is not initialized")
	}

	m := make(map[string]struct{})
	for i, entry := range d {
		w := strings.TrimSpace(entry)
		if w != entry {
			return fmt.Errorf("dictionary value %q has extra white space", entry)
		}
		if w == "" {
			return fmt.Errorf("dictionary value #%d is empty", i)
		}
		if _, ok := m[w]; ok {
			return fmt.Errorf("dictionary value %q is not unique", w)
		}
		m[w] = struct{}{}
	}
	return nil
}

// Load captures the first 256 words of a dictionary from an [io.Reader]. Lines starting with `//` are ignored.
func Load(r io.Reader) (d Dictionary, err error) {
	s := &scanner.Scanner{}
	s.Init(r)
	s.Error = func(s *scanner.Scanner, msg string) {
		err = errors.New(msg)
	}

	cursor := 0
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		word := strings.TrimSpace(s.TokenText())
		if strings.HasPrefix(word, "//") {
			continue // comment
		}
		if err != nil {
			return
		}
		d[cursor] = word
		cursor++
		if cursor > 255 {
			break
		}
	}
	return d, nil
}

func LoadFile(p string) (d Dictionary, err error) {
	handle, err := os.Open(p)
	if err != nil {
		return
	}
	defer handle.Close()
	return Load(handle)
}
