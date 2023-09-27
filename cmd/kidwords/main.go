package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/dkotik/kidwords"
	"github.com/spf13/pflag"
)

var (
	stdin   = pflag.Bool("stdin", false, "use data from OS standard in")
	intmod  = pflag.BoolP("integer", "i", false, "treat text as unsigned integer")
	reverse = pflag.BoolP("reverse", "r", false, "recover encoded data")
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
