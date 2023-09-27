package test

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func GoldenMust(t *testing.T, path string, value []byte) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if *update {
		_, err := io.Copy(f, bytes.NewReader(value))
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", path, err)
		}
		return // updated
	}

	expected, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", path, err)
	}

	if bytes.Compare(value, expected) != 0 {
		t.Log("Given:", string(value))
		t.Log("Expected:", string(expected))
		t.Log("Golden file:", path)
		t.Fatal("contents of the golden test data file does not match given value")
	}
}
