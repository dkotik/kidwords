package secret

import (
	"context"

	"github.com/dkotik/kidwords/service"
)

type Transaction interface {
	CommitOrRollback(context.Context, *error) error
}

type Repository interface {
	Create(context.Context, service.Secret) error
	Retrieve(context.Context, string) (service.Secret, error)
	Update(context.Context, service.Secret) error
	Delete(context.Context, string) error
	List(context.Context, string) ([]service.Secret, error)
	WithTransaction(context.Context, Transaction) (Repository, error)
}
