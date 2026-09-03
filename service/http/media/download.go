package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

//go:generate go run .

type DownloadTarget struct {
	FilePath string
	URL      string
}

var targets = []DownloadTarget{
	{
		FilePath: "htmx.min.js",
		URL:      "https://unpkg.com/htmx.org@latest/dist/htmx.min.js",
	},
	{
		FilePath: "bulma.min.css",
		URL:      "https://cdn.jsdelivr.net/npm/bulma@latest/css/bulma.min.css",
	},
}

func main() {
	for _, target := range targets {
		err := DownloadFile(target.FilePath, target.URL)
		if err != nil {
			fmt.Printf("failed to download %s: %v\n", target.FilePath, err)
			os.Exit(1)
		}
	}
}

func DownloadFile(filepath string, url string) error {
	// 1. Create the local file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// 2. Get the remote data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 3. Check for a valid server response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 4. Stream the data directly to the file without loading it all into memory
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file content: %w", err)
	}

	return nil
}
