package tgrid

import "io"

type Row []*Cell

func (r Row) Height() (height int) {
	for _, cell := range r {
		if height < cell.Height {
			height = cell.Height
		}
	}
	return
}

func (r Row) Write(w io.Writer) (n int, err error) {
	for i := 0; i < r.Height(); i++ {
		for j, cell := range r {
			if j == 0 {
				if _, err = w.Write([]byte("│")); err != nil {
					return
				}
			} else if _, err = w.Write([]byte("║")); err != nil {
				return
			}
			n++
			if err = cell.WriteLine(w); err != nil {
				return
			}
			n += cell.Width
		}
		if _, err = w.Write([]byte("│\n")); err != nil {
			return
		}
		n++
	}
	return
}
