package htservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/middleware/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var updateTitle = &i18n.LocalizeConfig{
	DefaultMessage: &i18n.Message{
		ID:    "updatePaperKeyTitle",
		Other: "Update Paper Key",
	},
}

type updateRequest struct {
	ID          string
	Name        string
	Description string
}

func (r *updateRequest) Validate(ctx context.Context) error {
	return nil
}

func NewUpdatePersonalSecretView(r SecretRepository) (http.Handler, error) {
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

			return &listResponse{
				Secrets: NewLabels(secrets),
				lc:      lc,
			}, nil
		},
		htadaptor.WithTemplate(templates.Lookup("update")),
	)
}
