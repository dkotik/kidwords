package http

import (
	"fmt"
	"net/http"

	"github.com/dkotik/htadaptor/staticfs"
)

//go:generate env GOOS=js GOARCH=wasm go build -o media/wkdw.v0.wasm ./wasm

// cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ./media

func NewWebAssemblyHandler(mux *http.ServeMux) error {
	wasmExecJS, err := assets.ReadFile("media/wasm_exec.js")
	if err != nil {
		return fmt.Errorf("unable to load wasm_exec source: %w", err)
	}
	mux.Handle("/wasm_exec.js", staticfs.NewFastFileSystemFileWithContentType(wasmExecJS, "text/javascript"))

	wasm, err := assets.ReadFile("media/wkdw.v0.wasm")
	if err != nil {
		return fmt.Errorf("unable to load wasm source: %w", err)
	}
	mux.Handle("/wkdw.v0.wasm", staticfs.NewFastFileSystemFileWithContentType(wasm, "application/wasm"))

	wasmHTML, err := assets.ReadFile("media/wasm.html")
	if err != nil {
		return fmt.Errorf("unable to load wasm_exec source: %w", err)
	}
	mux.Handle("/wasm.html", staticfs.NewFastFileSystemFileWithContentType(wasmHTML, "text/html"))

	return nil
}
