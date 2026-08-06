package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/dkotik/kidwords"
	"github.com/urfave/cli/v3"
)

var decode = &cli.Command{
	Name:      "decode",
	Usage:     "convert simple words into data",
	ArgsUsage: "\"-\" argument takes standard input",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		input := strings.Join(cmd.Args().Slice(), " ")
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
