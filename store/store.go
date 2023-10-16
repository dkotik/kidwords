/*
Package store provides interfaces and cryptographic primitives for [Store] implementations.
*/
package store

import (
	"context"
	"time"
)

type PaperKey struct {
	ID         string
	Owner      string
	Name       string
	SaltedHash string
	Created    time.Time
}

type Store interface {
	Create(ctx context.Context, p *PaperKey) error
	RetrieveAll(ctx context.Context, owner string) ([]*PaperKey, error)
	Delete(ctx context.Context, id string) error
}
