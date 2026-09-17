package secret

import (
	"context"
)

type Transaction interface {
	Close(*error)
}

type Repository interface {
	Create(context.Context, Secret) (string, error)
	Retrieve(context.Context, string) (Secret, error)
	Update(context.Context, Secret) error
	Delete(context.Context, string) error
	List(context.Context, string) ([]Secret, error)

	BeginTransaction(context.Context) (Repository, Transaction, error)
	WithTransaction(context.Context, Transaction) (Repository, error)
}

type User interface {
	GetID() string
}

type UserRepository interface {
	RetrieveUserByName(context.Context, string) (User, error)
	RetrieveUserByEmailAddress(context.Context, string) (User, error)

	BeginTransaction(context.Context) (Repository, Transaction, error)
	WithTransaction(context.Context, Transaction) (Repository, error)
}
