package kidwords

import (
	"fmt"
	"strconv"

	"github.com/dkotik/kidwords/dictionary"
)

type decoder struct {
	Nouns map[string]byte
	Verbs map[string]byte
}

type Decoder interface {
	Decode([]byte) ([]byte, error)
}

func NewDecoder(nouns, verbs dictionary.Dictionary) Decoder {
	return &decoder{
		Nouns: nouns.Reverse(),
		Verbs: verbs.Reverse(),
	}
}

func (d *decoder) buildShard(index int, nouns, verbs []byte) (_ []byte, err error) {
	if len(verbs) != 4 {
		return nil, fmt.Errorf("share %d: expected 4 verbs, got %d", index, len(verbs))
	}
	if !IsValid(nouns, verbs) {
		return nil, fmt.Errorf("share %d: invalid checksum", index)
	}
	return nouns, nil
}

func (d *decoder) Decode(data []byte) (b []byte, err error) {
	var (
		lastIndex int
		index     []byte
		word      []byte
		nouns     map[int][]byte
		verbs     map[int][]byte
		c         byte
		ok        bool
	)

	for _, c := range data {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if word != nil {
				c, ok = d.Nouns[string(word)]
				if ok {
					nouns[lastIndex] = append(nouns[lastIndex], c)
				} else {
					c, ok = d.Verbs[string(word)]
					if ok {
						verbs[lastIndex] = append(verbs[lastIndex], c)
					}
				}
				word = nil
			}
			index = append(index, c)
		default:
			if index != nil {
				lastIndex, err = strconv.Atoi(string(index))
				if err != nil {
					return nil, err
				}
				index = nil
			}
			word = append(word, c)
		}
	}

	c, ok = d.Nouns[string(word)]
	if ok {
		nouns[lastIndex] = append(nouns[lastIndex], c)
	} else {
		c, ok = d.Verbs[string(word)]
		if ok {
			verbs[lastIndex] = append(verbs[lastIndex], c)
		}
	}

	shares := make([]byte, 0, len(nouns))
	for index, ns := range nouns {
		vs, ok := verbs[index]
		if !ok {
			continue
		}
		share, err := d.buildShard(index, ns, vs)
		if err != nil {
			// return nil, err
			continue
		}
		shares = append(shares, share...)
	}
	return shares, nil
}
