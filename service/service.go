/*
Package service provides an authenticated portal for managing paper keys and regular passwords.
*/
package service

import (
	"context"
	"io"
)

type Component interface {
	Render(context.Context, io.Writer) error
}

type KeyValueRepository interface {
	Get(context.Context, []byte) ([]byte, error)
	Set(context.Context, []byte, []byte) error
	Update(context.Context, []byte, func([]byte) ([]byte, error)) error
	Delete(context.Context, []byte) error
}

type KeyKeyValueRepository interface {
	Get(context.Context, []byte, []byte) ([]byte, error)
	Set(context.Context, []byte, []byte, []byte) error
	Update(context.Context, []byte, []byte, func([]byte) ([]byte, error)) error
	Delete(context.Context, []byte, []byte) error
}

type domainAdaptor struct {
	kkv    KeyKeyValueRepository
	domain []byte
}

func NewKeyValueFromKeyKeyValueRepository(domain string, kkv KeyKeyValueRepository) KeyValueRepository {
	if len(domain) == 0 {
		panic("domain prefix is required")
	}
	if kkv == nil {
		panic("cannot use a <nil> key-key-value repository")
	}
	return &domainAdaptor{
		kkv:    kkv,
		domain: []byte(domain),
	}
}

func (a *domainAdaptor) Get(ctx context.Context, key []byte) ([]byte, error) {
	return a.kkv.Get(ctx, a.domain, key)
}

func (a *domainAdaptor) Set(ctx context.Context, key, value []byte) error {
	return a.kkv.Set(ctx, a.domain, key, value)
}

func (a *domainAdaptor) Update(ctx context.Context, key []byte, update func([]byte) ([]byte, error)) error {
	return a.kkv.Update(ctx, a.domain, key, update)
}

func (a *domainAdaptor) Delete(ctx context.Context, key []byte) error {
	return a.kkv.Delete(ctx, a.domain, key)
}

type Service struct {
	attempts AuthenticationAttemptRepository
	secrets  SecretRepository

	viewAuthenticationAttempts func([]AuthenticationAttempt) Component
	viewSecrets                func([]ArgonSecretLabel) Component
	viewCreateSecret           Component
	viewAuthenticate           Component
}
