/*
Package store provides interfaces and cryptographic primitives for [Store] implementations.
*/
package store

import (
	"context"
	"time"

	"log/slog"
)

type PaperKey struct {
	ID         string
	Owner      string
	Name       string
	SaltedHash string
	Created    time.Time
}

func (p *PaperKey) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", p.ID),
		slog.String("name", p.Name),
		slog.String("resource", "paperKey"),
	)
}

type Store interface {
	Create(ctx context.Context, p *PaperKey) error
	RetrieveAll(ctx context.Context, owner string) ([]*PaperKey, error)
	Delete(ctx context.Context, id string) error
	DeleteByOwner(ctx context.Context, owner string) error
}
