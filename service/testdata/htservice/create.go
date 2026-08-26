package htservice

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/middleware/session"
	"github.com/dkotik/kidwords"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var createTitle = &i18n.LocalizeConfig{
	DefaultMessage: &i18n.Message{
		ID:    "createPaperKeyTitle",
		Other: "Add Paper Key",
	},
}

type createResponse struct {
	Shards kidwords.Shards
	lc     *i18n.Localizer
}

func (r *createResponse) Head() (*htmlHead, error) {
	title, locale, err := r.lc.LocalizeWithTag(createTitle)
	if err != nil {
		return nil, err
	}
	base, _ := locale.Base()
	return &htmlHead{
		Title:  title,
		Locale: base.String(),
	}, nil
}

func NewCreatePersonalSecretView(
	r SecretRepository,
	keySize int,
	shares int,
	quorum int,
) (http.Handler, error) {
	if r == nil {
		return nil, errors.New("cannot use a <nil> secrets repository")
	}
	if keySize < 16 {
		return nil, errors.New("not enough key entropy: key is too short")
	}
	if shares < 1 {
		return nil, errors.New("not enough shares")
	}
	if quorum > shares {
		return nil, errors.New("quorum must be smaller than or the same as the number of shares")
	}

	return htadaptor.NewNullaryFuncAdaptor(
		func(ctx context.Context) (*createResponse, error) {
			id := session.UserID(ctx)
			if len(id) == 0 {
				return nil, session.ErrNoSessionInContext
			}
			lc, ok := htadaptor.LocalizerFromContext(ctx)
			if !ok {
				return nil, errors.New("no localizer in context")
			}
			secret := make([]byte, keySize)
			n, err := rand.Read(secret)
			if err != nil {
				return nil, err
			}
			if n != keySize {
				return nil, errors.New("not enough secure random bytes")
			}

			shards, err := kidwords.Split(
				string(secret),
				shares,
				quorum,
			)
			if err != nil {
				return nil, err
			}

			return &createResponse{
				Shards: shards,
				lc:     lc,
			}, nil
		},
		htadaptor.WithTemplate(templates.Lookup("create")),
	)
}
