package kidwords

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dkotik/kidwords/dictionary"
)

func TestEncode(t *testing.T) {
	tcs := []struct {
		Shards []Shard
	}{
		{
			Shards: []Shard{
				{Index: 0, Data: []byte("34534a23n"), Checksum: 5656},
			},
		},
		{
			Shards: []Shard{
				{Index: 0, Data: []byte("3453a5234"), Checksum: 5656},
				{Index: 1, Data: []byte("gho3452f4"), Checksum: 5656},
			},
		},
	}

	encoder, err := NewEncoder(
		dictionary.EnglishFourLetterNouns,
		dictionary.EnglishFourLetterVerbs,
		3,
	)
	if err != nil {
		t.Fatal(err)
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("tc%d", i), func(t *testing.T) {
			b := &bytes.Buffer{}
			_, _ = b.WriteString("\n")
			if err := encoder.Encode(b, tc.Shards); err != nil {
				t.Fatal(err)
			}
			t.Log(b.String())

			// t.Fail()
		})
	}
}
