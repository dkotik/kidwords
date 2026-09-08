package kidwords

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/dkotik/kidwords/dictionary"
)

func TestReader(t *testing.T) {
	r, err := NewReader(
		strings.NewReader(`idea...half...icon
      idea...crow...;grid!`),
		WithDictionary(&dictionary.EnglishFourLetterNouns),
	)
	if err != nil {
		t.Fatal(err)
	}
	b := &bytes.Buffer{}
	_, err = io.Copy(b, r)
	if err != nil {
		t.Fatal(err)
	}

	// goldie.New(t).Assert(t, "internal/testdata/readRaw.golden", b.Bytes())
}
