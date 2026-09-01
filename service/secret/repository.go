package secret

import (
	"context"
)

type Transaction interface {
	Close(*error)
}

type Repository interface {
	Create(context.Context, Secret) error
	Retrieve(context.Context, string) (Secret, error)
	Update(context.Context, Secret) error
	Delete(context.Context, string) error
	List(context.Context, string) ([]Secret, error)

	BeginTransaction(context.Context) (Repository, Transaction, error)
	WithTransaction(context.Context, Transaction) (Repository, error)
}
