package kidwords

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dkotik/kidwords/dictionary"
	"github.com/dkotik/kidwords/internal/shamir"
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

type tokenKind uint8

const (
	tokenBoundary tokenKind = iota
	tokenWord
	tokenIndex
)

func (d *decoder) Decode(data []byte) (b []byte, err error) {
	var (
		index       int
		token       []byte
		tokenKind   tokenKind
		nouns       = make(map[int][]byte)
		verbs       = make(map[int][]byte)
		shardErrors []error
	)

	consumeToken := func() {
		if len(token) == 0 {
			return
		}
		if tokenKind == tokenIndex {
			i, err := strconv.Atoi(string(token))
			if err != nil {
				shardErrors = append(shardErrors, fmt.Errorf("invalid index: %s", token))
			} else {
				index = i
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
			verbs[index] = append(verbs[index], c)
			token = nil
			return
		}
		shardErrors = append(shardErrors, fmt.Errorf("unknown word: %s", token))
	}

	for _, c := range data {
		switch c {
		case ' ', '\t', '\n', '.', ',', '\'', '"', '`', '|', '(', ')', '!', '?', '+', '-':
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
			token = append(token, c)
		default:
			switch tokenKind {
			case tokenWord:
			default:
				consumeToken()
				tokenKind = tokenWord
			}
			token = append(token, c)
		}
	}
	consumeToken() // last word

	shares := make([][]byte, 0, len(nouns))
	for index, ns := range nouns {
		vs, ok := verbs[index]
		if !ok {
			continue
		}
		share, err := d.buildShard(index, ns, vs)
		if err != nil {
			shardErrors = append(shardErrors, err)
			continue
		}
		shares = append(shares, share)
	}

	secret, err := shamir.Combine(shares)
	if err != nil {
		shardErrors = append(shardErrors, err)
		return nil, errors.Join(shardErrors...)
	}
	return secret, nil
}
