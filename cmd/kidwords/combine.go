package main

import (
	"bytes"
	"fmt"
	"strings"

	"io"
	"os"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/dictionary"
	"github.com/dkotik/kidwords/shamir"
	"github.com/urfave/cli/v2"
)

var combine = &cli.Command{
	Name:      "combine",
	Usage:     "recover the secret from a quorum of Shamir's Secret Sharing shards",
	ArgsUsage: "\"-\" argument takes standard input",
	Action: func(c *cli.Context) (err error) {
		input := strings.Join(c.Args().Slice(), " ")
		if input == "-" {
			b := &bytes.Buffer{}
			if _, err = io.Copy(b, os.Stdin); err != nil {
				return err
			}
			lines := bytes.Split(b.Bytes(), []byte("\n"))
			shards := make([][]byte, 0, len(lines))
			for _, line := range lines {
				b.Reset()
				r, err := kidwords.NewReader(bytes.NewReader(line))
				if err != nil {
					return err
				}
				if _, err = io.Copy(os.Stdout, r); err != nil {
					return err
				}
				if b.Len() > 0 {
					clone := make([]byte, b.Len())
					copy(clone, b.Bytes())
					shards = append(shards, clone)
				}
			}
			key, err := shamir.Combine(shards)
			if err != nil {
				return err
			}
			_, err = fmt.Printf("%s", key)
			return err
		}

		var shards [][]byte
		for {
			shard, more, err := scanShard(fmt.Sprintf("Collected %d shards", len(shards)))
			if err != nil {
				return err
			}
			shards = append(shards, shard)
			if !more {
				key, err := shamir.Combine(shards)
				if err != nil {
					return err
				}
				_, err = fmt.Printf("%s", key)
				return err
			}
		}
	},
}

func scanWord(prompt string) (string, error) {
	word, err := scanPassword(prompt)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(word)), nil
}

func scanShard(prompt string) (shard []byte, more bool, err error) {
	var words []string

top:
	for {
		word, err := scanWord(fmt.Sprintf("%s, %d words:", prompt, len(words)))
		if err != nil {
			return nil, false, err
		}
		for _, existing := range dictionary.EnglishFourLetterNouns {
			if existing == word {
				words = append(words, word)
				continue top
			}
		}
		switch word {
		case "":
			// fmt.Println(" ⚠ cannot use an empty word")
			fmt.Println(" ⚠ submit \"next\" to end the shard")
			fmt.Println(" ⚠ submit \"done\" to attempt recovery")
		case "next":
			shard, err := kidwords.ToBytes(strings.Join(words, " "))
			if err != nil {
				return nil, false, err
			}
			return shard, true, nil
		case "done":
			shard, err := kidwords.ToBytes(strings.Join(words, " "))
			if err != nil {
				return nil, false, err
			}
			return shard, false, nil
		default:
			fmt.Printf("word %q is not in the encoding dictionary\n", word)
		}
	}
}
