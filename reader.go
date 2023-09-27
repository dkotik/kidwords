package kidwords

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/dkotik/kidwords/dictionary"
)

func NewReader(r io.Reader, withOptions ...ReaderOption) (*Reader, error) {
	if r == nil {
		return nil, errors.New("cannot use a <nil> reader")
	}
	o := &readerOptions{}

	var err error
	for i, option := range withOptions {
		if err = option.applyReaderOption(o); err != nil {
			return nil, fmt.Errorf("cannot apply option %d to Kids Words reader: %w", i+1, err)
		}
	}

	if o.dictionary == nil {
		o.dictionary = (&dictionary.EnglishFourLetterNouns).Reverse()
	}

	return &Reader{
		r:          bufio.NewReader(r),
		dictionary: o.dictionary,
	}, nil
	// reader.Scanner.Error = func(s *scanner.Scanner, msg string) {
	// 	reader.scanErr = errors.New(msg)
	// }
	// reader.Scanner.IsIdentRune
	// reader.Scanner.Init(r)
	// return reader
}

type Reader struct {
	r *bufio.Reader
	// scanErr    error
	dictionary map[string]byte
}

func (r *Reader) Read(p []byte) (n int, err error) {
	var rn rune
	var b = strings.Builder{}

	for {
		rn, _, err = r.r.ReadRune()
		if err != nil {
			if err == io.EOF && b.Len() > 0 {
				data, ok := r.dictionary[b.String()]
				if !ok {
					return 0, fmt.Errorf("word %q is not in the dictionary", b.String())
				}
				p[0] = data
				return 1, nil
			}
			return 0, err
		}
		if !unicode.IsLetter(rn) {
			if b.Len() > 0 {
				data, ok := r.dictionary[b.String()]
				if !ok {
					return 0, fmt.Errorf("word %q is not in the dictionary", b.String())
				}
				p[0] = data
				return 1, nil
			}
			continue // skip non-letters
		}
		_, _ = b.WriteRune(rn)
	}
}

// func (r *Reader) Read(p []byte) (n int, err error) {
// 	// for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
// 	//
// 	// }
// 	for i := range p {
// 		c := r.Scan()
// 		if c == scanner.EOF {
// 			return n, io.EOF
// 		}
// 		if r.scanErr != nil {
// 			return n, r.scanErr
// 		}
// 		b, ok := r.dictionary[r.TokenText()]
// 		if !ok {
// 			return n, fmt.Errorf("word %q is not in chosen dictionary", r.TokenText())
// 		}
// 		p[i] = b
// 		n++
// 	}
// 	return
// }
