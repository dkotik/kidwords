package kidwords

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

type decoder struct {
	Nouns map[string]byte
	Verbs map[string]byte
}

type Decoder interface {
	Decode(string) ([]Shard, error)
}

func NewDecoder(nouns, verbs Dictionary) Decoder {
	return &decoder{
		Nouns: nouns.Reverse(),
		Verbs: verbs.Reverse(),
	}
}

type tokenKind uint8

const (
	tokenBoundary tokenKind = iota
	tokenWord
	tokenIndex
)

func (d *decoder) Decode(data string) (ss []Shard, err error) {
	var (
		n            int
		index        int
		highestIndex int
		token        []byte
		tokenKind    tokenKind
		nouns        = make(map[int][]byte)
		verbs        = make(map[int][]byte)
		shardErrors  []error
	)

	consumeToken := func() {
		if len(token) == 0 {
			return
		}
		if tokenKind == tokenIndex {
			i, err := strconv.Atoi(string(token))
			if err != nil {
				shardErrors = append(shardErrors, fmt.Errorf("invalid shard index: %s", token))
			} else if i < 1 || i > 255 {
				shardErrors = append(shardErrors, fmt.Errorf("shard index out of range: %d", index))
			} else {
				index = i
				if index > highestIndex {
					highestIndex = index
				}
			}
			token = nil
			return
		}
		c, ok := d.Nouns[string(token)]
		if ok {
			nouns[index] = append(nouns[index], c)
			token = nil
			return
		}
		c, ok = d.Verbs[string(token)]
		if ok {
			// fmt.Println(index, "verb", string(token))
			verbs[index] = append(verbs[index], c)
			token = nil
			return
		}
		shardErrors = append(shardErrors, fmt.Errorf("unknown word: %s", token))
	}

	runeBuf := make([]byte, 4)
	for _, c := range string(data) {
		switch c {
		case ' ', '\t', '\n', '·', '.', ',', '\'', '"', '`', '|', '(', ')', '!', '?', '+', '-', '[', ']', '{', '}', ':', ';', '\\', '/':
			switch tokenKind {
			case tokenBoundary: // do nothing
			default:
				consumeToken()
				tokenKind = tokenBoundary
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			switch tokenKind {
			case tokenIndex:
			default:
				consumeToken()
				tokenKind = tokenIndex
			}
			token = append(token, byte(c))
		default:
			switch tokenKind {
			case tokenWord:
			default:
				consumeToken()
				tokenKind = tokenWord
			}
			n = utf8.EncodeRune(runeBuf, c)
			token = append(token, runeBuf[:n]...)
		}
	}
	consumeToken() // last word
	if highestIndex == 0 {
		return nil, fmt.Errorf("no shards found")
	}

	ss = make([]Shard, 0, len(nouns))
	for index, ns := range nouns {
		vs, ok := verbs[index]
		if !ok {
			continue
		}
		shard, err := NewShard(uint8(index), ns, vs)
		if err != nil {
			shardErrors = append(shardErrors, err)
			continue
		}
		// fmt.Println(index, len(share))
		ss = append(ss, shard)
	}
	return ss, errors.Join(shardErrors...)
}
