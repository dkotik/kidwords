package kidwords

import (
	"bytes"
	"hash"
	"hash/crc32"
	"io"
)

var ChecksumTable = crc32.MakeTable(crc32.Koopman)

type checksumWriter struct {
	hash hash.Hash
	pass io.Writer
}

func (c *checksumWriter) Write(p []byte) (n int, err error) {
	n, err = c.hash.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}

	n, err = c.pass.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}
	return len(p), nil
}

func (c *checksumWriter) Close() (err error) {
	p := c.hash.Sum(nil)
	n, err := c.pass.Write(p)
	if err != nil {
		return
	}
	if n != c.hash.Size() {
		err = io.ErrShortWrite
		return
	}
	c.hash.Reset()
	return nil
}

func ChecksumWriter(w io.Writer) io.WriteCloser {
	return &checksumWriter{
		hash: crc32.New(ChecksumTable),
		pass: w,
	}
}

func ChecksumChop(b []byte) (remainder []byte, ok bool) {
	h := crc32.New(ChecksumTable)
	l := len(b) - h.Size()
	if l >= 0 {
		n, err := h.Write(b[:l])
		if err == nil && n == l {
			if bytes.Compare(h.Sum(nil), b[l:]) == 0 {
				return b[:l], true
			}
		}
	}
	return b, false
}
