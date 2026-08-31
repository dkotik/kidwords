package secret

import (
	"context"
)

type Transaction interface {
	CommitOrRollback(context.Context, *error) error
}

type Repository interface {
	Create(context.Context, Secret) error
	Retrieve(context.Context, string) (Secret, error)
	Update(context.Context, Secret) error
	Delete(context.Context, string) error
	List(context.Context, string) ([]Secret, error)

	BeginTransaction(context.Context) (Transaction, error)
	WithTransaction(context.Context, Transaction) (Repository, error)
}
