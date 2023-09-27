package test

import (
	"fmt"
	"testing"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/shamir"
)

func TestShardCombination(t *testing.T) {
	t.Skip("dictionary outdated")
	shards := [][]byte{
		[]byte("boat neck mail risk baby milk seal bird silk team"),
		[]byte("step fire view clay city baby hall wall icon base"),
		[]byte("hero cold firm luck idea deck army card boat soil"),
		[]byte("rock dust tape corn heat stop aunt corn plan cold"),
		[]byte("neck pike drum boot bush joke land duck fear mask"),
		[]byte("note gust rope stew tank iron army foil hint golf"),
		[]byte("junk skin form neck trap goat neck junk cell core"),
		[]byte("shot pool bell sage deck time moon loan link past"),
	}

	for i, shard := range shards {
		data, err := kidwords.ToBytes(string(shard))
		if err != nil {
			t.Fatal(err)
		}
		shards[i] = data
	}

	for i := 0; i <= 5; i++ {
		fmt.Printf("%d: %x\n", i+1, shards[i])
		fmt.Printf("%d: %x\n", i+2, shards[i+1])
		fmt.Printf("%d: %x\n", i+3, shards[i+2])

		data, err := shamir.Combine(shards[i : i+3])
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}

	t.Fatal("check")
}
