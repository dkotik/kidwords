package kidwords

import (
	"math"

	"github.com/dkotik/kidwords/internal/shamir"
)

type Secret []Shard

func NewSecret(data []byte, parts, threshold int) (Secret, error) {
	shares, err := shamir.Split(data, parts, threshold)
	if err != nil {
		return nil, err
	}
	secret := make(Secret, len(shares))
	for i, share := range shares {
		secret[i] = Shard{
			Index:    uint8(i),
			Data:     share,
			Checksum: NewChecksum(share),
		}
		// fmt.Printf("shard %d: %d\n", i, share[len(share)-1])
	}
	return secret, nil
}

func (s Secret) GetFingerprint() []byte {
	var (
		shard Shard
		fp    = make([]byte, len(s))
		i, j  = 0, 0
	)

	for i, shard = range s {
		j = len(shard.Data) - 1
		if j == -1 {
			continue
		}
		fp[i] = shard.Data[j]
	}
	return fp
}

func (s Secret) MatchFingerprint(fp []byte) bool {
	i := 0
	l := uint8(max(len(fp), math.MaxUint8) - 1)
	for _, shard := range s {
		if shard.Index > l {
			return false
		}
		i = len(shard.Data) - 1
		if i == -1 {
			return false
		}
		if shard.Data[i] != fp[shard.Index] {
			return false
		}
	}
	return true
}
