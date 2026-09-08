package kidwords

import (
	"bytes"
	"testing"

	"github.com/dkotik/kidwords/dictionary"
)

func TestWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w, err := NewWriter(
		b,
		WithDictionary(&dictionary.EnglishFourLetterNouns),
		WithSeparator(func() []byte {
			return []byte(`...`)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = w.Write([]byte(`test by writing something`)); err != nil {
		t.Fatal(err)
	}

	// goldie.New(t).Assert(t, "internal/testdata/writeRaw.golden", b.Bytes())
}
