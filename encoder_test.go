package kidwords

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	tcs := []struct {
		Shards []Shard
	}{
		{
			Shards: []Shard{
				{Index: 0, Data: []byte("34534a23n"), Checksum: []byte{0x02, 0x38, 0x02, 0x38}},
			},
		},
		{
			Shards: []Shard{
				{Index: 0, Data: []byte("3453a5234"), Checksum: []byte{0x02, 0x38, 0x02, 0x38}},
				{Index: 1, Data: []byte("gho3452f4"), Checksum: []byte{0x02, 0x38, 0x02, 0x38}},
			},
		},
	}

	encoder, err := NewEncoder(
		EnglishFourLetterNouns,
		EnglishFourLetterVerbs,
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
		})
	}
}

func TestForEmptyRows(t *testing.T) {
	tcs := []struct {
		Count  int
		Quorum int
		Secret []byte
	}{
		{Count: 12, Quorum: 6, Secret: []byte("34534a23n")},
		{Count: 12, Quorum: 6, Secret: []byte("3453a5234a")},
		{Count: 12, Quorum: 6, Secret: []byte("gho3452f4bn")},
		{Count: 12, Quorum: 6, Secret: []byte("gho3452f4bn1")},
	}
	encoder, err := NewEncoder(
		EnglishFourLetterNouns,
		EnglishFourLetterVerbs,
		3,
	)
	if err != nil {
		t.Fatal(err)
	}

	for i, tc := range tcs {
		shards, err := NewSecret(tc.Secret, tc.Count, tc.Quorum)
		if err != nil {
			t.Fatal(err)
		}

		table := encoder.MakeTable(shards)
		for _, row := range table {
			for _, cell := range row {
				empty := true
				for _, word := range cell.Words {
					if word != "" && word != blank {
						empty = false
						break
					}
				}
				t.Logf("row: %2d %v", cell.Index, cell.Words)
				if empty {
					t.Fatalf("test case %d: found empty row in table for secret %q", i+1, tc.Secret)
				}
			}
		}
	}
}
