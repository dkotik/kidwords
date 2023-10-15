//go:build generate

package main

import (
	"io"
	"os"
	"regexp"
)

func main() {
	handle, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}
	defer handle.Close()
	modFile, err := io.ReadAll(handle)
	if err != nil {
		panic(err)
	}

	m := regexp.MustCompile(`github.com/dkotik/kidwords (v\d+\.\d+\.\d+)`).FindSubmatch(modFile)

	w, err := os.Create("version.gen.go")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if _, err = w.Write([]byte("package main\n\nconst version = \"")); err != nil {
		panic(err)
	}
	if _, err = w.Write(m[1]); err != nil {
		panic(err)
	}
	if _, err = w.Write([]byte("\"\n")); err != nil {
		panic(err)
	}
}
