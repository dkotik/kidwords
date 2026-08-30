package service

import "errors"

type options struct {
	Localizer Localizer
	KeyLength uint8
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
