/*
Package service authenticates accounts with paper keys and manages account paper keys.

It is designed to integrate securely with an existing HTTP service.
*/
package service

import (
	"context"
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

type Repository interface {
	Push(context.Context, Secret) error
	Pull(context.Context, string, string) ([]Secret, error)
	Delete(context.Context, string) error
}
