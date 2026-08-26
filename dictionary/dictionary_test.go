package dictionary

import "testing"

func TestValidateDictionaries(t *testing.T) {
	known := make(map[string]string)

	for _, word := range EnglishFourLetterNouns {
		if dictionary, ok := known[word]; ok {
			t.Fatalf("word %q is not unique (found in %s)", word, dictionary)
		}
		known[word] = "nouns"
	}

	for _, word := range EnglishFourLetterVerbs {
		if dictionary, ok := known[word]; ok {
			t.Fatalf("word %q is not unique (found in %s)", word, dictionary)
		}
		known[word] = "verbs"
	}

	var err error
	if err = EnglishFourLetterNouns.Validate(); err != nil {
		t.Fatal("English four letter nouns contain a flaw:", err)
	}
	if err = EnglishFourLetterVerbs.Validate(); err != nil {
		t.Fatal("English four letter verbs contain a flaw:", err)
	}
}
