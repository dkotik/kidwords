//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/dkotik/kidwords"
)

// bytesFromJS converts a JavaScript string, Array, or typed array to bytes.
func bytesFromJS(value js.Value) ([]byte, error) {
	if value.Type() == js.TypeString {
		return []byte(value.String()), nil
	}

	if value.IsNull() || value.IsUndefined() {
		return nil, fmt.Errorf("expected a string or byte array")
	}

	length := value.Get("length")
	if length.Type() != js.TypeNumber || length.Int() < 0 {
		return nil, fmt.Errorf("expected a string or byte array")
	}

	data := make([]byte, length.Int())
	for i := range data {
		item := value.Index(i)
		if item.Type() != js.TypeNumber || item.Int() < 0 || item.Int() > 255 {
			return nil, fmt.Errorf("byte at index %d is out of range", i)
		}
		data[i] = byte(item.Int())
	}
	return data, nil
}

func bytesToJS(data []byte) js.Value {
	result := js.Global().Get("Uint8Array").New(len(data))
	for i, value := range data {
		result.SetIndex(i, int(value))
	}
	return result
}

// split creates parts shares from a secret. It is exposed as goSplit(secret,
// parts, threshold), where secret may be a string or a byte array. The result
// is an array of Uint8Array values.
func split(this js.Value, args []js.Value) any {
	if len(args) != 3 {
		return "Error: secret, parts, and threshold are required"
	}
	secret, err := bytesFromJS(args[0])
	if err != nil {
		return "Error: " + err.Error()
	}
	if args[1].Type() != js.TypeNumber || args[2].Type() != js.TypeNumber {
		return "Error: parts and threshold must be numbers"
	}

	secretParts, err := kidwords.NewSecret(secret, args[1].Int(), args[2].Int())
	if err != nil {
		return "Error: " + err.Error()
	}

	encoder, err := kidwords.NewEncoder(kidwords.EnglishFourLetterNouns, kidwords.EnglishFourLetterVerbs, 3)
	if err != nil {
		return "Error: " + err.Error()
	}

	b := &bytes.Buffer{}
	if err = encoder.Encode(b, secretParts); err != nil {
		return "Error: " + err.Error()
	}

	return b.String()
}

// combine reconstructs a secret from an array of byte arrays. It is exposed
// as goCombine(shards) and returns a Uint8Array.
func combine(this js.Value, args []js.Value) any {
	if len(args) != 1 || args[0].IsNull() || args[0].IsUndefined() {
		return "Error: a string of shards is required"
	}

	decoder := kidwords.NewDecoder(
		kidwords.EnglishFourLetterNouns,
		kidwords.EnglishFourLetterVerbs,
	)
	shards, _, ok := decoder.Decode(args[0].String())
	if !ok {
		return "Error: not enough shards"
	}

	secret, err := kidwords.Combine(shards)
	if err != nil {
		return "Error: " + err.Error()
	}
	// result := js.Global().Get("Array").New(len(secretParts))
	// for i, shard := range secretParts {
	// 	result.SetIndex(i, bytesToJS(shard.Data))
	// }
	return string(secret)
}

// This function is kept as a small example of calling Go from JavaScript.
func add(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return "Error: Two numbers required"
	}

	return args[0].Int() + args[1].Int()
}

func main() {
	js.Global().Set("goAdd", js.FuncOf(add))
	js.Global().Set("splitSecret", js.FuncOf(split))
	js.Global().Set("combineShards", js.FuncOf(combine))

	// Keep the Go program alive so JavaScript can keep calling the functions.
	select {}
}
