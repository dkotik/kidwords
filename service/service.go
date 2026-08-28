/*
Package service authenticates accounts with paper keys and manages account paper keys.

It is designed to integrate securely with an existing HTTP service.
*/
package service

import (
	"context"
	"log/slog"
	"time"
)

type Secret struct {
	ID        string
	UserID    string
	Name      string
	Type      string
	Hash      Argon2Hash
	CreatedAt time.Time
	UsedAt    time.Time
}

func (p *Secret) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", p.ID),
		slog.String("name", p.Name),
		slog.String("type", p.Type),
	)
}

type Repository interface {
	Create(context.Context, Secret) error
	Retrieve(context.Context, string, string) ([]Secret, error)
	Update(context.Context, Secret) error
	Delete(context.Context, string) error
}
