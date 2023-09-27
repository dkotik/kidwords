package kidwords

import (
	"fmt"
	"io"

	"github.com/dkotik/kidwords/dictionary"
)

type Writer struct {
	io.Writer
	separator  SeparatorFunc
	dictionary *dictionary.Dictionary
}

func NewWriter(out io.Writer, withOptions ...WriterOption) (*Writer, error) {
	if out == nil {
		out = io.Discard
	}
	o := &writerOptions{}

	var err error
	for i, option := range withOptions {
		if err = option.applyWriterOption(o); err != nil {
			return nil, fmt.Errorf("cannot apply option %d to Kids Words writer: %w", i+1, err)
		}
	}

	if o.dictionary == nil {
		o.dictionary = &dictionary.EnglishFourLetterNouns
	}
	if o.separator == nil {
		o.separator = func() []byte {
			return []byte(" ")
		}
	}

	return &Writer{
		Writer:     out,
		separator:  o.separator,
		dictionary: o.dictionary,
	}, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var (
		j, l int
		sep  []byte
	)
	for _, c := range p {
		sep = w.separator()
		l = len(sep)
		if l > 0 {
			j, err = w.Writer.Write(sep)
			if err != nil {
				return
			}
			if j != l {
				return n, io.ErrShortWrite
			}
		}

		j, err = w.Writer.Write([]byte(w.dictionary[c]))
		if err != nil {
			return
		}
		if j != len(w.dictionary[c]) {
			return n, io.ErrShortWrite
		}
		n++
	}
	return n, nil
}

// func NewWriter(w io.Writer) io.WriteCloser {
// 	return ChecksumWriter(&Writer{
// 		Writer: w,
// 		separator: func(wordCount int) []byte {
// 			if wordCount%4 == 0 {
// 				if wordCount == 0 {
// 					return nil
// 				}
// 				return []byte("\n")
// 			}
// 			return []byte(" ")
// 		},
// 		dictionary: defaultDictionary,
// 	})
// }

// func WriteInt(w io.Writer, n int64) (err error) {
// 	if n < 0 {
// 		return errors.New("negative integers cannot be reliabily encoded to bytes")
// 	}
//
// 	x := new(big.Int).SetInt64(n)
// 	wrapped := NewWriter(w)
// 	_, err = wrapped.Write(x.Bytes())
// 	if err != nil {
// 		return
// 	}
// 	return wrapped.Close()
// }

// func SeparatorHTML(perRow int) SeparatorFunc {
// 	return func(wordCount int) []byte {
// 		if wordCount%perRow == 0 {
// 			if wordCount == 0 {
// 				return nil
// 			}
// 			return []byte("<br />")
// 		}
// 		return []byte("&nbsp;")
// 	}
// }
