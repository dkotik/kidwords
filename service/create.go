package service

import (
	"context"
	"crypto/rand"
	"strings"
	"time"
	"uuid"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FormCreateKey struct {
	lc           *i18n.Localizer
	User         User
	Secret       kidwords.Secret
	KeyName      string
	KeyNameError string
}

func (s *Service) createKeyFormView(ctx context.Context) (any, error) {
	user, lc, err := s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	return &FormCreateKey{
		lc:   lc,
		User: user,
	}, nil
}

func (s *Service) createKeyFormPost(ctx context.Context, name string) (_ *FormCreateKey, err error) {
	form := &FormCreateKey{}
	form.User, form.lc, err = s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	form.KeyName = strings.TrimSpace(name)
	if form.KeyName == "" {
		form.KeyNameError, err = form.lc.LocalizeMessage(&i18n.Message{
			ID:    "KidwordsErrorKeyNameRequired",
			Other: "key name is empty",
		})
		return form, err
	}

	secretBytes := make([]byte, s.keyLength)
	if _, err = rand.Read(secretBytes); err != nil {
		form.KeyNameError = err.Error()
		return form, err
	}
	argonHash, err := secret.NewArgonHash(secretBytes)
	if err != nil {
		return form, err
	}

	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return form, err
	}
	defer tx.Close(&err)
	rp, err = s.repository.WithTransaction(ctx, tx)
	if err != nil {
		return form, err
	}
	keys, err := rp.List(ctx, form.User.GetID())
	if err != nil {
		return form, err
	}
	if len(keys) >= s.keyCountLimit {
		form.KeyNameError, err = form.lc.LocalizeMessage(&i18n.Message{
			ID:    "KidwordsErrorTooManyKeys",
			Other: "too many paper keys, please delete one",
		})
		return form, err
	}
	err = rp.Create(ctx, secret.Secret{
		ID:         uuid.New().String(),
		UserID:     form.User.GetID(),
		Name:       form.KeyName,
		Type:       PaperKeyType,
		SaltedHash: argonHash.String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return form, err
	}

	form.Secret, err = kidwords.NewSecret(
		secretBytes,
		s.shardCount,
		s.quorumCount,
	)
	if err != nil {
		return form, err
	}

	return form, nil
}
