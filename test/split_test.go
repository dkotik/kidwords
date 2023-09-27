package test

import (
	"bytes"
	"testing"

	"github.com/dkotik/kidwords/shamir"
)

func TestSplit(t *testing.T) {
	secret := []byte(`a12345678`)
	shards, err := shamir.Split(secret, 8, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i, shard := range shards {
		t.Logf("%d: %x", i+1, shard)
	}

	for i := 0; i < 5; i++ {
		data, err := shamir.Combine(shards[i : i+3])
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(secret, data) != 0 {
			t.Fatal("recovered data does not match")
		}
	}
	// t.Fatal("checking")
}
