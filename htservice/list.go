package htservice

import (
	"context"
	"errors"
	"net/http"
	"sort"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/middleware/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var listTitle = &i18n.LocalizeConfig{
	DefaultMessage: &i18n.Message{
		ID:    "listPaperKeysTitle",
		Other: "Paper Keys",
	},
}

type listResponse struct {
	Secrets []ArgonSecretLabel
	lc      *i18n.Localizer
}

func (l *listResponse) Head() (*htmlHead, error) {
	title, locale, err := l.lc.LocalizeWithTag(listTitle)
	if err != nil {
		return nil, err
	}
	base, _ := locale.Base()
	return &htmlHead{
		Title:  title,
		Locale: base.String(),
	}, nil
}

func NewPersonalSecretsView(r SecretRepository) (http.Handler, error) {
	if r == nil {
		return nil, errors.New("cannot use a <nil> secrets repository")
	}

	return htadaptor.NewNullaryFuncAdaptor(
		func(ctx context.Context) (*listResponse, error) {
			id := session.UserID(ctx)
			if len(id) == 0 {
				return nil, session.ErrNoSessionInContext
			}
			lc, ok := htadaptor.LocalizerFromContext(ctx)
			if !ok {
				return nil, errors.New("no localizer in context")
			}

			secrets, err := r.ListSecrets(ctx, id)
			if err != nil {
				return nil, err
			}
			labels := make([]ArgonSecretLabel, len(secrets))
			for i, secret := range secrets {
				labels[i] = secret.Label()
			}
			sort.Slice(labels, func(a, b int) bool {
				return labels[a].Created.Before(labels[b].Created)
			})

			return &listResponse{
				Secrets: labels,
				lc:      lc,
			}, nil
		},
	)
}
