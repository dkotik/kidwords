package main

import (
	"bytes"
	"io"
)

func compress(b []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	compressor, err := zstd.NewWriter(
		b,
		zstd.WithEncoderLevel(zstd.Best),
		zstd.WithEncoderCRC(false),
	)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(compressor, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	if err = compressor.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
