package kidwords

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/dkotik/kidwords/dictionary"
	"github.com/dkotik/kidwords/internal"
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

	internal.GoldenMustMatch(t, "internal/testdata/readRaw.golden", b.Bytes())
}
