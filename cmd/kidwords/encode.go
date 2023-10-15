package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/dkotik/kidwords"
	"github.com/urfave/cli/v2"
)

var encode = &cli.Command{
	Name:      "encode",
	Usage:     "convert input into simple words",
	ArgsUsage: "\"-\" argument takes standard input",
	Action: func(c *cli.Context) error {
		w, err := kidwords.NewWriter(os.Stdout)
		if err != nil {
			return err
		}
		if strings.Join(c.Args().Slice(), " ") == "-" {
			_, err = io.Copy(w, os.Stdin)
			return err
		}

		secret, err := scanPassword("Enter secret: ")
		if err != nil {
			return err
		}
		if _, err = os.Stdout.Write([]byte("Encoded:")); err != nil {
			return err
		}
		if _, err = io.Copy(w, bytes.NewReader(secret)); err != nil {
			return err
		}
		_, err = os.Stdout.Write([]byte("\n"))
		return err
	},
}
