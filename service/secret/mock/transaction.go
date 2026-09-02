package mock

import (
	"context"

	"github.com/dkotik/kidwords/service/secret"
)

type tx struct{}

func (tx tx) Close(*error) {}

func (m *mock) BeginTransaction(ctx context.Context) (secret.Repository, secret.Transaction, error) {
	return m, tx{}, nil
}

func (m *mock) WithTransaction(ctx context.Context, tx secret.Transaction) (secret.Repository, error) {
	return m, nil
}
