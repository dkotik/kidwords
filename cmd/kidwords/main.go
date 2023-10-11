package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/shamir"
	"github.com/dkotik/kidwords/tgrid"
	"github.com/spf13/pflag"
)

var (
	stdin   = pflag.Bool("stdin", false, "use data from session standard input")
	intmod  = pflag.BoolP("integer", "i", false, "treat text as unsigned integer")
	reverse = pflag.BoolP("reverse", "r", false, "recover encoded data")
	quorum  = pflag.UintP("quorum", "q", 0, "split output into Shamir Secret Sharing shards")
	help    = pflag.BoolP("help", "h", false, "display help message")
)

func output(b []byte) {
	fmt.Println(kidwords.FromBytes(b))
}

func translate(s string) {
	b, err := kidwords.ToBytes(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	pflag.Parse()

	if *help {
		fmt.Fprint(os.Stderr, "kidwords: Paper key encoder and decoder.\n\n  kidwords [WORD1] [WORD2] ...\n\n")
		pflag.PrintDefaults()
		return
	}

	words := pflag.Args()
	if len(words) > 0 {
		if quorum != nil {
			input := strings.Join(words, " ")
			// shards, err := shamirSplit([]byte(input), uint8(*quorum))
			shards, err := shamir.Split([]byte(input), 12, int(*quorum))
			if err != nil {
				panic(err)
			}

			i := 0
			grid, err := tgrid.NewGrid(4, 3, func() (*tgrid.Cell, error) {
				words, err := kidwords.FromBytes(shards[i])
				if err != nil {
					return nil, err
				}
				i++
				return tgrid.NewCellFromBytes([]byte(words), 18), nil
			})
			if err != nil {
				panic(err)
			}

			fmt.Printf(" 🔑 Pick any %d shards:\n", *quorum)
			if _, err = grid.Write(os.Stdout); err != nil {
				panic(err)
			}

			// for i, shard := range shards {
			// 	compressed := []byte(shard)
			// 	words, err := kidwords.FromBytes(compressed)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	fmt.Printf("#%d: %d\n", i+1, int(compressed[len(compressed)-1]))
			// 	fmt.Printf("#%d: %s\n", i+1, words)
			// }
			fmt.Printf("Versioned key recovery command: go run github.com/dkotik/kidwords/cmd/kidwords@%s recover\n", commit)
			// output, err := kidwords.FromString(input)
			return
		}

		if *reverse {
			translate(strings.Join(words, " "))
			return
		}

		for _, take := range words {
			output([]byte(take))
		}
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	in := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			in <- scanner.Text()
		}
	}()

	for {
		select {
		case <-c:
			return
		case value := <-in:
			if *reverse {
				translate(value)
				continue
			}
			output([]byte(value))
		}
	}
}
