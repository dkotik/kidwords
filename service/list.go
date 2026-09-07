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

func (v keyListView) Description() (string, error) {
	return v.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsListDescription",
			Other: "Paper keys are the best account recovery mechanism. Create a key. Print the key and store it in a safe place. If you lose access to your account in the future, you can sign in using any paper key.",
		},
	})
}

func (v keyListView) CreateKeyTitle() (string, error) {
	return newFormCreateKeyTitle(v.lc)
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
