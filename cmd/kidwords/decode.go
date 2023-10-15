package main

import (
	"io"
	"os"
	"strings"

	"github.com/dkotik/kidwords"
	"github.com/urfave/cli/v2"
)

var decode = &cli.Command{
	Name:      "decode",
	Usage:     "convert simple words into data",
	ArgsUsage: "\"-\" argument takes standard input",
	Action: func(c *cli.Context) error {
		input := strings.Join(c.Args().Slice(), " ")
		if input == "-" {
			r, err := kidwords.NewReader(os.Stdin)
			if err != nil {
				return err
			}
			_, err = io.Copy(os.Stdout, r)
			return err
		}
		r, err := kidwords.NewReader(strings.NewReader(input))
		if err != nil {
			return err
		}
		_, err = io.Copy(os.Stdout, r)
		return err
	},
}
