package kidwords

import (
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
