/*
Package service authenticates accounts with paper keys and manages account paper keys.

It is designed to integrate securely with an existing HTTP service.
*/
package service

import (
	"context"
	"errors"

	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type User interface {
	GetID() string
	GetName() string
}

type Authenticator interface {
	Authenticate(context.Context) (User, error)
}

type Localizer interface {
	GetLocalizer(context.Context) *i18n.Localizer
}

type LocalizerFunc func(context.Context) *i18n.Localizer

func (f LocalizerFunc) GetLocalizer(ctx context.Context) *i18n.Localizer {
	return f(ctx)
}

type Service struct {
	authenticator Authenticator
	repository    secret.Repository
	localizer     Localizer
	keyLength     int
}

func New(
	authenticator Authenticator,
	repository secret.Repository,
	withOptions ...Option,
) (*Service, error) {
	o := &options{}
	for _, opt := range withOptions {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if o.Localizer == nil {
		lc := i18n.NewLocalizer(i18n.NewBundle(language.English))
		o.Localizer = LocalizerFunc(func(context.Context) *i18n.Localizer {
			return lc
		})
	}
	if o.KeyLength == 0 {
		o.KeyLength = 24
	}
	return &Service{
		authenticator: authenticator,
		repository:    repository,
		localizer:     o.Localizer,
		keyLength:     int(o.KeyLength),
	}, nil
}

func (s *Service) unpackContext(ctx context.Context) (u User, l *i18n.Localizer, err error) {
	l = s.localizer.GetLocalizer(ctx)
	if l == nil {
		return nil, nil, errors.New("no localizer in request context")
	}
	u, err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return u, l, nil
}
