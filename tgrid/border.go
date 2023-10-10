package tgrid

import (
	"bytes"
	"io"
)

func (g Grid) WriteTopBorder(w io.Writer) (n int, err error) {
	if len(g) < 1 {
		return
	}
	for i, cell := range g[0] {
		if i == 0 {
			if _, err = w.Write([]byte("┌")); err != nil {
				return
			}
			n++
		} else {
			if _, err = w.Write([]byte("╥")); err != nil {
				return
			}
			n++
		}
		for j := 0; j < cell.Width; j++ {
			if _, err = w.Write([]byte("─")); err != nil {
				return
			}
			n++
		}
	}
	if _, err = w.Write([]byte("┐\n")); err != nil {
		return
	}
	n += 2
	return
}

func (g Grid) MiddleBorder() []byte {
	if len(g) < 1 {
		return nil
	}
	b := &bytes.Buffer{}
	for i, cell := range g[0] {
		if i == 0 {
			_, _ = b.WriteRune('╞')
		} else {
			_, _ = b.WriteRune('╬')
		}
		for j := 0; j < cell.Width; j++ {
			_, _ = b.WriteRune('═')
		}
	}
	_, _ = b.WriteRune('╡')
	_, _ = b.WriteRune('\n')
	return b.Bytes()
}

func (g Grid) WriteBottomBorder(w io.Writer) (n int, err error) {
	if len(g) < 1 {
		return
	}
	for i, cell := range g[0] {
		if i == 0 {
			if _, err = w.Write([]byte("└")); err != nil {
				return
			}
			n++
		} else {
			if _, err = w.Write([]byte("╨")); err != nil {
				return
			}
			n++
		}
		for j := 0; j < cell.Width; j++ {
			if _, err = w.Write([]byte("─")); err != nil {
				return
			}
			n++
		}
	}
	if _, err = w.Write([]byte("┘\n")); err != nil {
		return
	}
	n += 2
	return
}
