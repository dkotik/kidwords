/*
Package main is a command line utility for encoding durable and accessible paper keys.
*/
package main

import (
	"context"
	"fmt"
	"os"

	"runtime/debug"

	"github.com/urfave/cli/v3"
)

var commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return "latest"
}()

func main() {
	if err := (&cli.Command{
		Name:                  "kidwords",
		Usage:                 "durable and accessible paper key codec\n<https://github.com/dkotik/kidwords>",
		Version:               version(),
		HideHelp:              false,
		HideVersion:           false,
		EnableShellCompletion: true,
		Suggest:               true,
		Commands: []*cli.Command{
			split,
			combine,
			encode,
			decode,
		},
	}).Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error: %s.\n", err.Error())
		os.Exit(1)
	}
}

// func main() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, os.Kill)
// 	in := make(chan string)
// 	go func() {
// 		scanner := bufio.NewScanner(os.Stdin)
// 		for scanner.Scan() {
// 			in <- scanner.Text()
// 		}
// 	}()
//
// 	for {
// 		select {
// 		case <-c:
// 			return
// 		case value := <-in:
// 			if *reverse {
// 				translate(value)
// 				continue
// 			}
// 			output([]byte(value))
// 		}
// 	}
// }

func version() string {
	v := "dev"
	if info, ok := debug.ReadBuildInfo(); ok {
		v = info.Main.Version
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				v = v + "-" + setting.Value
				break
			}
		}
	}
	return v
}
