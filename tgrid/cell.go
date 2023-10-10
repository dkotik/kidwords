package tgrid

import (
	"bytes"
	"io"
)

type Line []byte

type CellGenerator func() (*Cell, error)

type Cell struct {
	Next   int
	Width  int
	Height int
	Lines  []Line
}

func NewCell(fromLines ...Line) *Cell {
	cell := &Cell{
		Height: len(fromLines),
		Lines:  fromLines,
	}

	length := 0
	for _, line := range fromLines {
		if length = len(line); length > cell.Width {
			cell.Width = length
		}
	}
	return cell
}

func NewCellFromBytes(b []byte, lineWidth int) *Cell {
	wrapped := WordWrap(b, lineWidth)
	split := bytes.Split(wrapped, []byte("\n"))
	lines := make([]Line, len(split))
	for i, line := range split {
		lines[i] = Line(line)
	}
	return NewCell(lines...)
}

func (c *Cell) WriteLine(w io.Writer) (err error) {
	// log.Println(c.Next, c.Height)
	if c.Next >= c.Height {
		if err = c.WriteFiller(w, c.Width); err != nil {
			return
		}
		return
	}
	n, err := io.Copy(w, bytes.NewReader(c.Lines[c.Next]))
	if err != nil {
		return err
	}
	if err = c.WriteFiller(w, c.Width-int(n)); err != nil {
		return err
	}
	c.Next++
	return err
}

func (c *Cell) WriteFiller(w io.Writer, n int) (err error) {
	for ; n > 0; n-- {
		if _, err = w.Write([]byte(" ")); err != nil {
			return
		}
	}
	return
}
