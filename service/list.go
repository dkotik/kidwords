package service

import (
	"context"

	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type keyListView struct {
	lc        *i18n.Localizer
	Locale    string
	User      User
	Title     string
	PaperKeys []secret.Secret
}

func (s *Service) List(ctx context.Context) (view keyListView, err error) {
	view.User, view.lc, err = s.unpackContext(ctx)
	if err != nil {
		return view, err
	}
	var tag language.Tag
	view.Title, tag, err = view.lc.LocalizeWithTag(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsListTitle",
			Other: "My Paper Keys",
		},
	})
	if err != nil {
		return view, err
	}
	view.Locale = tag.String()

	view.PaperKeys, err = s.repository.List(ctx, view.User.GetID())
	if err != nil {
		return view, err
	}
	// for _, key := range view.PaperKeys {
	// 	if key.UserID == "" {
	// 		return view, errors.New("empty userID")
	// 	}
	// }

	return view, nil
}
