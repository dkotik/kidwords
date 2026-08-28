package kidwords

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/dkotik/kidwords/dictionary"
)

const blank = "    "

type Cell struct {
	Index uint8
	Words [4]string
}

type Table [][]Cell

type Encoder interface {
	MakeTable(shards []Shard) Table
	Encode(w io.Writer, shards []Shard) error
}

type encoder struct {
	Nouns   dictionary.Dictionary
	Verbs   dictionary.Dictionary
	Columns int
}

func NewEncoder(nouns, verbs dictionary.Dictionary, columns int) (_ Encoder, err error) {
	if err = nouns.Validate(); err != nil {
		return nil, err
	}
	if err = verbs.Validate(); err != nil {
		return nil, err
	}
	if columns < 1 {
		return nil, errors.New("at least one column is required")
	}

	return &encoder{
		Nouns:   nouns,
		Verbs:   verbs,
		Columns: columns,
	}, nil
}

func (e *encoder) shardToWords(shard Shard) (words []string) {
	count := len(shard.Data)
	index := 0
	var c byte
	words = make([]string, 0, count+4)
	verbBytes := int32ToBytesBigEndian(shard.Checksum)
	verbsCount := 4

	for _, c = range shard.Data {
		words = append(words, e.Nouns[c])
		index++
		if index%3 == 0 && verbsCount > 0 && index > 0 {
			words = append(words, e.Verbs[verbBytes[4-verbsCount]])
			// words = append(words, "****")
			verbsCount--
			// index += 4
			// index++
		}
	}

	// remaining verbs
	for ; verbsCount > 0; verbsCount-- {
		words = append(words, e.Verbs[verbBytes[4-verbsCount]])
		// words = append(words, "????")
	}

	// remaining blanks
	for range 4 - (count % 4) {
		words = append(words, blank)
	}

	return words
}

func (e *encoder) MakeTable(shards []Shard) (table Table) {
	var (
		shard Shard
		words []string
		cells []Cell
		// rows  [][]Cell
		// cell  Cell
		i, j, index int
	)
	for chunk := range slices.Chunk(shards, e.Columns) {
		shard = chunk[0]
		for words = range slices.Chunk(e.shardToWords(shard), 4) {
			cells = make([]Cell, e.Columns)
			cells[0] = Cell{
				Index: shard.Index + 1,
				Words: [4]string{
					words[0],
					words[1],
					words[2],
					words[3],
				},
			}
			table = append(table, cells)
		}

		for i, shard = range chunk[1:] {
			j = index
			i++
			for words = range slices.Chunk(e.shardToWords(shard), 4) {
				table[j][i] = Cell{
					Index: shard.Index + 1,
					Words: [4]string{
						words[0],
						words[1],
						words[2],
						words[3],
					},
				}
				j++
			}
		}
		index++
	}
	return
}

func (e *encoder) Encode(w io.Writer, shards []Shard) (err error) {
	padding := len(shards) / 10
	table := e.MakeTable(shards)
	for _, row := range table {
		for i, cell := range row {
			if cell.Index == 0 {
				_, err = w.Write(bytes.Repeat([]byte("·"), padding+1+5+5+5+5))
				if err != nil {
					return err
				}
				continue
			}

			_, err = w.Write(bytes.Repeat([]byte(" "), padding-(int(cell.Index)-1)/10))
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w, "%d", cell.Index)
			if err != nil {
				return err
			}
			for _, word := range cell.Words {
				_, err = w.Write([]byte(" "))
				if err != nil {
					return err
				}
				_, err = w.Write([]byte(word))
				if err != nil {
					return err
				}
			}
			if i < e.Columns {
				_, err = w.Write([]byte(" "))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func int32ToBytesBigEndian(val int32) []byte {
	return []byte{
		byte(val >> 24),
		byte(val >> 16),
		byte(val >> 8),
		byte(val),
	}
}
