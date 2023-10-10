package tgrid

import (
	"os"
	"testing"
)

func TestGeneratedGrid(t *testing.T) {
	grid, err := NewGrid(4, 3, func() (*Cell, error) {
		line := Line([]byte("1 2 3"))
		return NewCell(line, line, line), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err = grid.Write(os.Stdout); err != nil {
		t.Fatal(err)
	}
	// t.Fatal("============================================")
}

func TestManualGrid(t *testing.T) {
	line := Line([]byte("1 2 3"))
	line2 := Line([]byte("1 2 3 4 5"))

	row1 := Row([]*Cell{
		NewCell(line, line, line),
		NewCell(line, line, line),
		// NewCell(line, line, line),
	})
	row2 := Row([]*Cell{
		NewCell(line),
		NewCell(line, line, line2),
		NewCell(line, line, line2),
	})
	row3 := Row([]*Cell{
		NewCell(line2, line, line),
		NewCell(line, line, line),
		// NewCell(line, line, line),
	})

	grid := Grid([]Row{row1, row2, row3})
	grid.Normalize()

	_, err := grid.Write(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	// t.Fatal("============================================")
}
