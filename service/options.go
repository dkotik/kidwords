package service

import (
	"errors"
)

const (
	DefaultKeyCount    = 12
	DefaultKeyLength   = 12
	DefaultShardCount  = 12
	DefaultQuorumCount = 6
)

type options struct {
	Localizer     Localizer
	KeyCountLimit uint8
	KeyLength     uint8
	ShardCount    uint8
	QuorumCount   uint8
}

type Option func(*options) error

func WithLocalizer(localizer Localizer) Option {
	return func(o *options) error {
		if localizer == nil {
			return errors.New("nil localizer")
		}
		if o.Localizer != nil {
			return errors.New("localizer already set")
		}
		o.Localizer = localizer
		return nil
	}
}

func WithKeyLength(keyLength uint8) Option {
	return func(o *options) error {
		if keyLength < 1 {
			return errors.New("key length must be positive")
		}
		if o.KeyLength != 0 {
			return errors.New("key length already set")
		}
		o.KeyLength = keyLength
		return nil
	}
}

func WithKeyLimit(limit uint8) Option {
	return func(o *options) error {
		if limit < 1 {
			return errors.New("key count must be positive")
		}
		if o.KeyCountLimit != 0 {
			return errors.New("key count already set")
		}
		o.KeyCountLimit = limit
		return nil
	}
}

func WithShardCount(total, quorum uint8) Option {
	return func(o *options) error {
		if quorum < 1 {
			return errors.New("quorum count must be positive")
		}
		if total < 1 {
			return errors.New("shard count must be positive")
		}
		if total < quorum {
			return errors.New("shard count must be greater than quorum count")
		}

		if o.ShardCount != 0 {
			return errors.New("shard count already set")
		}
		if o.QuorumCount != 0 {
			return errors.New("quorum count already set")
		}
		o.ShardCount = total
		o.QuorumCount = quorum
		return nil
	}
}
