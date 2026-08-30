package service

import (
	"context"
	"crypto/rand"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FormCreateKey struct {
	lc           *i18n.Localizer
	User         User
	KeyName      string
	KeyNameError string
}

func (f *FormCreateKey) Validate(ctx context.Context) error {
	return nil
}

func (s *Service) createKeyFormView(ctx context.Context, name string) (any, error) {
	user, lc, err := s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	return &FormCreateKey{
		lc:      lc,
		User:    user,
		KeyName: name,
	}, nil
}

func (s *Service) createKeyFormPost(ctx context.Context, form *FormCreateKey) (_ *FormCreateKey, err error) {
	form.User, form.lc, err = s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	form.KeyName = strings.TrimSpace(form.KeyName)
	if form.KeyName == "" {
		form.KeyNameError, err = form.lc.LocalizeMessage(&i18n.Message{
			ID:    "KidwordsErrorKeyNameRequired",
			Other: "key name is empty",
		})
		return form, err
	}

	b := make([]byte, s.keyLength)
	if _, err = rand.Read(b); err != nil {
		form.KeyNameError = err.Error()
		return form, err
	}

	return form, nil
}
