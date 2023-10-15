package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dkotik/kidwords"
	"github.com/urfave/cli/v2"
)

var split = &cli.Command{
	Name:      "split",
	Usage:     "split input into Shamir's Secret Sharing shards",
	ArgsUsage: "\"-\" argument takes standard input",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "shards",
			Aliases: []string{"s"},
			Usage:   "the number of shards to create",
			Value:   12,
			Action: func(ctx *cli.Context, n int) error {
				if n < 2 || n > 256 {
					return fmt.Errorf("Flag shards value %d out of range[2-256]", n)
				}
				return nil
			},
		},
		&cli.IntFlag{
			Name:    "quorum",
			Aliases: []string{"q"},
			Usage:   "the number of shards required to recover the secret",
			Value:   4,
			Action: func(ctx *cli.Context, n int) error {
				if n < 2 || n > 256 {
					return fmt.Errorf("Flag quorum value %d out of range[2-256]", n)
				}
				return nil
			},
		},
		&cli.IntFlag{
			Name:    "columns",
			Aliases: []string{"c"},
			Usage:   "the number of table columns in the output grid",
			Value:   3,
			Action: func(ctx *cli.Context, n int) error {
				if n < 1 || n > 12 {
					return fmt.Errorf("Flag columns value %d out of range[1-12]", n)
				}
				return nil
			},
		},
		&cli.IntFlag{
			Name:    "wrap",
			Aliases: []string{"w"},
			Usage:   "maximum shard line length",
			Value:   18,
			Action: func(ctx *cli.Context, n int) error {
				if n < 4 || n > 128 {
					return fmt.Errorf("Flag wrap value %d out of range[4-128]", n)
				}
				return nil
			},
		},
	},
	Action: func(c *cli.Context) error {
		input := strings.Join(c.Args().Slice(), " ")
		if input == "-" {

		}

		parts := c.Value("shards").(int)
		threshold := c.Value("quorum").(int)
		shards, err := kidwords.Split(input, parts, threshold)
		if err != nil {
			return err
		}
		if _, err = fmt.Printf(" 🔑 Pick any %d shards:\n", threshold); err != nil {
			return err
		}

		columns := c.Value("columns").(int)
		wrap := c.Value("wrap").(int)
		if _, err = shards.Grid(columns, wrap).Write(os.Stdout); err != nil {
			return err
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
		_, err = fmt.Printf("go run github.com/dkotik/kidwords/cmd/kidwords@%s combine\n", commit)
		return err
	},
}
