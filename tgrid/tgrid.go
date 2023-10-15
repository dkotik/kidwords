/*
Package tgrid represents tables as simple ASCII art. Example:

	┌──╥──╥──╥───┐
	│12║12║12║123│
	╞══╬══╬══╬═══╡
	│12║12║12║123│
	╞══╬══╬══╬═══╡
	│12║12║12║123│
	└──╨──╨──╨───┘
*/
package tgrid

import (
	"bytes"
	"io"
)

// Grid is an ASCII rendering of table data.
type Grid []Row

func NewGrid(columns, rows int, from CellGenerator) (Grid, error) {
	g := Grid(make([]Row, rows))
	for i := 0; i < rows; i++ {
		row := make([]*Cell, columns)
		for j := 0; j < columns; j++ {
			cell, err := from()
			if err != nil {
				return nil, err
			}
			row[j] = cell
		}
		g[i] = row
	}
	g.Normalize()
	return g, nil
}

func (g Grid) Normalize() {
	if len(g) == 0 {
		return
	}

	mostCells, currentCells, addWidths := 0, 0, 0
	columnWidths := make([]int, len(g[0]))
	for _, row := range g {
		if currentCells = len(row); mostCells < currentCells {
			mostCells = currentCells
		}

		// expand columnWidths, if neccesary
		if addWidths = currentCells - len(columnWidths); addWidths > 0 {
			columnWidths = append(
				columnWidths,
				make([]int, addWidths)...,
			)
		}

		for i, cell := range row {
			if addWidths = columnWidths[i]; addWidths < cell.Width {
				columnWidths[i] = cell.Width
			}
		}
	}

	for i, row := range g {
		// grow column widths
		for j, cell := range row {
			cell.Width = columnWidths[j]
		}

		// grow additional cells
		for currentCells = len(row); currentCells < mostCells; currentCells++ {
			g[i] = append(g[i], &Cell{
				Width: columnWidths[currentCells],
			})
		}
	}
}

func (g Grid) Write(w io.Writer) (n int, err error) {
	border := g.MiddleBorder()

	added := 0
	added64 := int64(0)

	added, err = g.WriteTopBorder(w)
	n += added
	if err != nil {
		return
	}

	lastRow := len(g) - 1
	for i, row := range g {
		added, err = row.Write(w)
		n += added
		if err != nil {
			return
		}

		if i == lastRow {
			// write
		} else {
			added64, err = io.Copy(w, bytes.NewReader(border))
			n += int(added64)
			if err != nil {
				return
			}
		}
	}

	added, err = g.WriteBottomBorder(w)
	n += added
	if err != nil {
		return
	}
	return
}
