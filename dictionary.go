package kidwords

import (
	"errors"
	"fmt"
	"strings"
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
