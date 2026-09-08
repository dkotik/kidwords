package service

import (
	"context"
	"time"

	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type keyView struct {
	ID             string
	UserID         string
	Name           string
	Type           string
	Fingerprint    []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastAcceptedAt time.Time
}

type keyListView struct {
	lc        *i18n.Localizer
	Locale    string
	User      User
	Title     string
	PaperKeys []keyView
}

func newKeyView(s secret.Secret) keyView {
	return keyView{
		ID:             s.ID,
		UserID:         s.UserID,
		Name:           s.Name,
		Type:           s.Type,
		Fingerprint:    s.Fingerprint,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		LastAcceptedAt: s.LastAcceptedAt,
	}
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
	return v.lc.LocalizeMessage(createButtonLabel)
}

func (f keyListView) DeleteButtonLabel() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: deleteButtonLabel,
	})
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

	secrets, err := s.repository.List(ctx, view.User.GetID())
	if err != nil {
		return view, err
	}
	view.PaperKeys = make([]keyView, len(secrets))
	for i, s := range secrets {
		view.PaperKeys[i] = newKeyView(s)
	}
	return view, nil
}
