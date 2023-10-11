package kidwords

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	shards, err := Split("somethingElse", 12, 8)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = shards.Grid(6, 18).Write(os.Stdout); err != nil {
		t.Fatal(err)
	}
	if err = shards.WriteHTML(os.Stdout, 3); err != nil {
		t.Fatal(err)
	}
	// t.Fatal("show")
}

func compress(r io.Reader) ([]byte, error) {
	b := &bytes.Buffer{}
	zr, err := gzip.NewWriterLevel(b, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(zr, r); err != nil {
		return nil, err
	}
	if err = zr.Flush(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func compressionRatio(originalLen, compressedLen int) float32 {
	return float32(compressedLen) / float32(originalLen)
}

func TestProveCompressionIsPointless(t *testing.T) {
	random := "jaksdASIUY3298sakhj*@&SkjASDjkhndsj"
	compressed, err := compress(strings.NewReader(random))
	if err != nil {
		t.Fatal(err)
	}
	if r := compressionRatio(len(random), len(compressed)); r < .9 {
		t.Fatalf("compression ratio turned out to be favorable: %.2f", r)
	} else {
		t.Logf("unfavorable compression ratio: %.2f", r)
	}

	bip39key := "goddess return math panther sustain black fatigue tortoise vast steel fiction scare"
	compressed, err = compress(strings.NewReader(bip39key))
	if err != nil {
		t.Fatal(err)
	}
	if r := compressionRatio(len(bip39key), len(compressed)); r < .9 {
		t.Fatalf("compression ratio turned out to be favorable: %.2f", r)
	} else {
		t.Logf("unfavorable compression ratio: %.2f", r)
	}

	bip39key = "swear expand tourist debris drink found word adjust skin input hawk extend"
	compressed, err = compress(strings.NewReader(bip39key))
	if err != nil {
		t.Fatal(err)
	}
	if r := compressionRatio(len(bip39key), len(compressed)); r < .9 {
		t.Fatalf("compression ratio turned out to be favorable: %.2f", r)
	} else {
		t.Logf("unfavorable compression ratio: %.2f", r)
	}
	// t.Fatal("compressed poorly")
}
