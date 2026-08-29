package kidwords

import (
	"bytes"
	"testing"

	"github.com/dkotik/kidwords/dictionary"
)

func TestDecoder(t *testing.T) {
	secretBytes := []byte("secret")
	secret, err := NewSecret(secretBytes, 12, 6)
	if err != nil {
		t.Fatal(err)
	}
	if len(secret) != 12 {
		t.Fatal("secret length is not 12")
	}

	encoder, err := NewEncoder(dictionary.EnglishFourLetterNouns, dictionary.EnglishFourLetterVerbs, 3)
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	if err = encoder.Encode(b, secret); err != nil {
		t.Fatal(err)
	}
	if b.Len() == 0 {
		t.Fatal("empty buffer")
	}

	decoder := NewDecoder(
		dictionary.EnglishFourLetterNouns,
		dictionary.EnglishFourLetterVerbs,
	)

	decoded, err := decoder.Decode(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if string(decoded) != string(secretBytes) {
		t.Log("original:", string(secretBytes))
		t.Log("decoded:", string(decoded))
		t.Fatal("decoded secret does not match the original")
	}
}
