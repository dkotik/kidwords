package main

import (
	"github.com/extism/go-pdk"
)

// https://github.com/extism/go-pdk
//
// go install github.com/extism/cli/extism@latest
// brew tap tinygo-org/tools
// brew install tinygo
// tinygo build -o plugin.wasm -target wasi main.go
// extism call plugin.wasm greet --input "Benjamin" --wasi

//export greet
func greet() int32 {
	input := pdk.Input()
	greeting := `Hello, ` + string(input) + `!`
	pdk.OutputString(greeting)
	return 0
}

func main() {}
