package kidwords

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

type InvalidShardIndexError struct {
	Token string
}

func (e InvalidShardIndexError) Error() string {
	return fmt.Sprintf("invalid shard index: %s", e.Token)
}

type UnknownWordError struct {
	Token string
}

func (e UnknownWordError) Error() string {
	return fmt.Sprintf("word is neither in the noun nor in the verb dictionary: %s", e.Token)
}

type decoder struct {
	Nouns map[string]byte
	Verbs map[string]byte
}

type Decoder interface {
	Decode(string) ([]Shard, []error, bool)
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

func (d *decoder) Decode(data string) (ss []Shard, errs []error, ok bool) {
	var (
		n            int
		index        int
		highestIndex int
		c            byte
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
			if err != nil || (i < 1 || i > 255) {
				errs = append(errs, InvalidShardIndexError{
					Token: string(token),
				})
			} else {
				index = i
				if index > highestIndex {
					highestIndex = index
				}
			}
			goto clear
		}
		c, ok = d.Nouns[string(token)]
		if ok {
			nouns[index] = append(nouns[index], c)
			goto clear
		}
		c, ok = d.Verbs[string(token)]
		if ok {
			// fmt.Println(index, "verb", string(token))
			verbs[index] = append(verbs[index], c)
			goto clear
		}
		errs = append(errs, UnknownWordError{
			Token: string(token),
		})

	clear:
		token = nil
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
		return ss, errs, false
	}

	ss = make([]Shard, 0, len(nouns))
	for index, ns := range nouns {
		token, ok = verbs[index]
		if !ok {
			continue
		}
		shard, err := NewShard(uint8(index), ns, token)
		if err != nil {
			shardErrors = append(shardErrors, err)
			continue
		}
		ss = append(ss, shard)
	}
	return ss, errs, true
}
