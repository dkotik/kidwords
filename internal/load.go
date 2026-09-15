package internal

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
)

// Load captures the first 256 words of a dictionary from an [io.Reader]. Lines starting with `//` are ignored.
func Load(r io.Reader) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		s := bufio.NewScanner(r)

		found := 0
		for s.Scan() && found < 256 {
			word := strings.TrimSpace(s.Text())
			if strings.HasPrefix(word, "//") {
				continue // comment
			}
			if len(word) != 4 {
				yield("", fmt.Errorf("word %q is %d characters long instead of 4", word, len(word)))
				return
			}
			if !yield(word, nil) {
				return
			}
			found++
		}

		if err := s.Err(); err != nil {
			yield("", err)
			return
		}

		if found < 255 {
			yield("", fmt.Errorf("not enough dictionary values: %d vs 255", found))
			return
		}
	}
}

func LoadFile(p string) iter.Seq2[string, error] {
	handle, err := os.Open(p)
	if err != nil {
		return func(yield func(string, error) bool) {
			yield("", err)
		}
	}
	defer handle.Close()
	return Load(handle)
}
