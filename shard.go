package kidwords

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/dkotik/kidwords/internal/shamir"
	"github.com/dkotik/kidwords/tgrid"
)

type Shard struct {
	Index    uint8
	Data     []byte
	Checksum []byte
}

func NewShard(index uint8, nouns, verbs []byte) (s Shard, err error) {
	if len(verbs) != 4 {
		return s, fmt.Errorf("share %d: expected 4 verbs, got %d", index, len(verbs))
	}
	if !IsValid(nouns, verbs) {
		return s, fmt.Errorf("share %d: invalid checksum", index)
	}
	s.Index = index
	s.Data = nouns
	s.Checksum = verbs
	return s, nil
}

func Combine(ss []Shard) ([]byte, error) {
	shares := make([][]byte, len(ss))
	for i, shard := range ss {
		shares[i] = shard.Data
	}
	return shamir.Combine(shares)
}

type Shards []string

func (s Shards) Grid(columns, wrap int) tgrid.Grid {
	total := len(s)
	rows := int(math.Ceil(float64(total) / float64(columns)))
	g := tgrid.Grid(make([]tgrid.Row, rows))
	for i := 0; i < rows; i++ {
		row := make([]*tgrid.Cell, columns)
		for j := 0; j < columns; j++ {
			index := i + j
			if index >= total {
				g.Normalize()
				return g
			}
			cell := tgrid.NewCellFromBytes([]byte(s[i+j]), wrap)
			// cell, err := from()
			// if err != nil {
			// 	return nil, err
			// }
			row[j] = cell
		}
		g[i] = row
	}
	g.Normalize()
	return g
}

func (s Shards) WriteHTML(w io.Writer, columns int) (err error) {
	if _, err = w.Write([]byte("<table>")); err != nil {
		return err
	}
	if _, err = w.Write([]byte("<tr>")); err != nil {
		return err
	}
	for i, shard := range s {
		if i > 0 && i%columns == 0 {
			if _, err = w.Write([]byte("</tr><tr>")); err != nil {
				return err
			}
		}
		if _, err = w.Write([]byte("<td>")); err != nil {
			return err
		}
		if _, err = io.Copy(w, strings.NewReader(shard)); err != nil {
			return err
		}
		if _, err = w.Write([]byte("</td>")); err != nil {
			return err
		}
	}
	if _, err = w.Write([]byte("</tr>")); err != nil {
		return err
	}
	_, err = w.Write([]byte("</table>"))
	return err
}

func (s Shards) Write(w io.Writer) (int, error) {
	return s.Grid(4, 18).Write(w)
}

func (s Shards) String() string {
	b := &bytes.Buffer{}
	if _, err := s.Write(b); err != nil {
		return "<unserializable>"
	}
	return b.String()
}

// func Split(
// 	key string,
// 	total,
// 	quorum int,
// 	withOptions ...WriterOption,
// ) (shards Shards, err error) {
// 	raw, err := shamir.Split([]byte(key), total, quorum)
// 	if err != nil {
// 		return nil, err
// 	}
// 	shards = make([]string, len(raw))

// 	for i, shard := range raw {
// 		encoded, err := FromBytes(shard, withOptions...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		shards[i] = encoded
// 	}

// 	return shards, nil
// }
