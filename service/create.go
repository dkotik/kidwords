package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"uuid"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FormCreateKey struct {
	lc           *i18n.Localizer
	Locale       string
	User         User
	Secret       kidwords.Secret
	KeyName      string
	KeyNameLabel string
	KeyNameError string
}

func (f *FormCreateKey) Title() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsCreateFormTitle",
			Other: "Create New Paper Key",
		},
	})
}

func (s *Service) newFormCreateKey(ctx context.Context) (f *FormCreateKey, err error) {
	f = &FormCreateKey{}
	f.User, f.lc, err = s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	knl, tag, err := f.lc.LocalizeWithTag(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsKeyNameLabel",
			Other: "key name",
		},
	})
	f.KeyNameLabel = knl
	f.Locale = tag.String()
	return f, nil
}

func (s *Service) createKeyFormView(ctx context.Context) (any, error) {
	f, err := s.newFormCreateKey(ctx)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) CreateKeyFormPost(ctx context.Context, name string) (_ *FormCreateKey, err error) {
	form, err := s.newFormCreateKey(ctx)
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
		return form, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	userID := form.User.GetID()
	keys, err := rp.List(ctx, userID)
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
	_, err = rp.Create(ctx, secret.Secret{
		ID:         uuid.New().String(),
		UserID:     userID,
		Name:       form.KeyName,
		Type:       secret.TypePaperKey,
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
