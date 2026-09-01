package service

import (
	"context"

	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type keyListView struct {
	lc        *i18n.Localizer
	User      User
	PaperKeys []secret.Secret
}

func (s *Service) list(ctx context.Context) (view keyListView, err error) {
	view.User, view.lc, err = s.unpackContext(ctx)
	if err != nil {
		return view, err
	}

	view.PaperKeys, err = s.repository.List(ctx, view.User.GetID())
	if err != nil {
		return view, err
	}

	return view, nil
}
