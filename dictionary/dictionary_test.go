package dictionary

import "testing"

func TestValidateDictionaries(t *testing.T) {
	var err error

	if err = EnglishFourLetterNouns.Validate(); err != nil {
		t.Fatal("English four letter nouns contain a flaw:", err)
	}
}
