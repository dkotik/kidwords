//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

// This function will be callable from JavaScript
func add(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return "Error: Two numbers required"
	}

	// Pull values out of the JavaScript arguments array
	num1 := args[0].Int()
	num2 := args[1].Int()

	return num1 + num2
}

func main() {
	fmt.Println("Hello from Go WebAssembly!")

	// Expose the Go function to the global JavaScript scope (window object)
	js.Global().Set("goAdd", js.FuncOf(add))

	// Keep the Go program alive so JavaScript can keep calling the functions
	select {}
}
