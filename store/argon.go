package store

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	DefaultArgonKeyLength  = 128
	DefaultArgonTimeCost   = 1
	DefaultArgonMemoryCost = 64 * 1024 // recommended by x/crypto/argon2
	DefaultArgonThreads    = 4
)

// ArgonHash is a parameterized salted hash used for storing keys.
type ArgonHash struct {
	Type            string
	Version         uint8
	TimeCost        uint32
	MemoryCost      uint32
	ParallelThreads uint8
	Salt            []byte
	Secret          []byte
}

// NewArgonHash creates an Argon2id hash using default parameters.
func NewArgonHash(key []byte) (*ArgonHash, error) {
	return NewCustomArgonHash(
		key,
		DefaultArgonTimeCost,
		DefaultArgonMemoryCost,
		DefaultArgonThreads,
	)
}

// NewCustomArgonHash creates an Argon hash.
func NewCustomArgonHash(
	key []byte,
	timeCost uint32,
	memoryCost uint32,
	parallelThreads uint8,
) (*ArgonHash, error) {
	salt := &bytes.Buffer{}
	if _, err := io.CopyN(salt, rand.Reader, DefaultArgonKeyLength); err != nil {
		return nil, err
	}
	return &ArgonHash{
		Type:            "argon2id",
		Version:         19,
		TimeCost:        timeCost,
		MemoryCost:      memoryCost,
		ParallelThreads: parallelThreads,
		Salt:            salt.Bytes(),
		Secret: argon2.Key(
			key,
			salt.Bytes(),
			timeCost,
			memoryCost,
			parallelThreads,
			DefaultArgonKeyLength),
	}, nil
}

// Match hashes the given key using [ArgonHash] parameters and compares the result with [ArgonHash.Secret].
func (a *ArgonHash) Match(key []byte) (bool, error) {
	switch a.Type {
	case "argon2d":
		hash := argon2.Key(
			key,
			a.Salt,
			a.TimeCost,
			a.MemoryCost,
			a.ParallelThreads,
			DefaultArgonKeyLength)
		return bytes.Compare(hash, a.Secret) == 0, nil
	case "argon2id":
		hash := argon2.IDKey(
			key,
			a.Salt,
			a.TimeCost,
			a.MemoryCost,
			a.ParallelThreads,
			DefaultArgonKeyLength)
		return bytes.Compare(hash, a.Secret) == 0, nil
	default:
		return false, fmt.Errorf("hash type %q is not supported", a.Type)
	}
}

// String serializes the [ArgonHash] using format `$<type>$v=<version>$m=<memory>,t=<time>,p=<parallel>$<salt>$<secret>`.
func (a *ArgonHash) String() string {
	return fmt.Sprintf(`$%s$v=%d$m=%d,t=%d,p=%d$%s$%s`,
		a.Type,
		a.Version,
		a.MemoryCost,
		a.TimeCost,
		a.ParallelThreads,
		base64.RawStdEncoding.EncodeToString(a.Salt),
		base64.RawStdEncoding.EncodeToString(a.Secret),
	)
}

// ParseArgonHash constructs an [ArgonHash] from a serialized string following the format `$<type>$v=<version>$m=<memory>,t=<time>,p=<parallel>$<salt>$<secret>`.
func ParseArgonHash(h string) (result *ArgonHash, err error) {
	fragments := strings.FieldsFunc(h, func(r rune) bool {
		return r == '$'
	})
	if l := len(fragments); l != 5 {
		return nil, fmt.Errorf("hash contains %d sections instead of 5", l)
	}

	if !strings.HasPrefix(fragments[1], "v=") {
		return nil, errors.New("hash version is corrupt")
	}
	version, err := strconv.ParseUint(fragments[1][2:], 10, 64)
	if err != nil || version == 0 {
		return nil, errors.New("hash version cannot be parsed")
	}

	result = &ArgonHash{
		Type:    fragments[0],
		Version: uint8(version),
	}

	result.Salt, err = base64.RawStdEncoding.DecodeString(fragments[3])
	if err != nil {
		return nil, err
	}
	result.Secret, err = base64.RawStdEncoding.DecodeString(fragments[4])
	if err != nil {
		return nil, err
	}

	for _, pair := range strings.FieldsFunc(fragments[2], func(r rune) bool {
		return r == ','
	}) {
		key, value, _ := strings.Cut(pair, "=")
		switch key { // m=<memory>,t=<time>,p=<parallel>
		case "m":
			cost, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot parse memory cost: %w", err)
			}
			result.MemoryCost = uint32(cost)
		case "t":
			cost, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot parse time cost: %w", err)
			}
			result.TimeCost = uint32(cost)
		case "p":
			cost, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot parse parallel threads: %w", err)
			}
			result.ParallelThreads = uint8(cost)
		}
	}
	return result, nil
}
