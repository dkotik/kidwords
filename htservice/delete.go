package htservice

import (
	"context"
	"errors"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/extract"
	"github.com/dkotik/htadaptor/middleware/session"
)

func NewDeletePersonalSecretView(r SecretRepository) (http.Handler, error) {
	if r == nil {
		return nil, errors.New("cannot use a <nil> secrets repository")
	}
	extractor, err := extract.NewQueryValueExtractor("delete")
	if err != nil {
		return nil, err
	}

	return htadaptor.NewUnaryStringFuncAdaptor(
		func(ctx context.Context, secretID string) (*listResponse, error) {
			id := session.UserID(ctx)
			if len(id) == 0 {
				return nil, session.ErrNoSessionInContext
			}
			lc, ok := htadaptor.LocalizerFromContext(ctx)
			if !ok {
				return nil, errors.New("no localizer in context")
			}
			err := r.DeleteSecret(ctx, id, secretID)
			if err != nil {
				return nil, err
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
		extractor,
	)
}
