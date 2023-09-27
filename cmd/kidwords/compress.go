package main

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
)

func compress(in []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	compressor, err := zstd.NewWriter(
		b,
		zstd.WithEncoderLevel(zstd.SpeedBestCompression),
		zstd.WithEncoderCRC(false),
	)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(compressor, bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	if err = compressor.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
