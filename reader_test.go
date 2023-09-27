package kidwords

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/dkotik/kidwords/dictionary"
	"github.com/dkotik/kidwords/test"
)

func TestReader(t *testing.T) {
	r, err := NewReader(
		strings.NewReader(`idea...half...icon
      idea...cell...;grid!`),
		WithDictionary(&dictionary.EnglishFourLetterNouns),
	)
	b := &bytes.Buffer{}
	_, err = io.Copy(b, r)
	if err != nil {
		t.Fatal(err)
	}

	test.GoldenMust(t, "test/testdata/readRaw.golden", b.Bytes())
}
