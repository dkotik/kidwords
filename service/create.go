package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var createButtonLabel = &i18n.Message{
	ID:    "KidwordsCreateFormTitle",
	Other: "Create New Paper Key",
}

type FormCreateKey struct {
	lc     *i18n.Localizer
	Title  string
	Name   string
	Locale string
	User   User
	Secret kidwords.Table
	Error  string
	Quorum int
}

func (f *FormCreateKey) Description() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsCreateFormDescription",
			Other: "Print and store this key in a safe place. The key is split into validatable shards. You will need at least {{ .Quorum }} valid shards to gain access to your account.",
		},
		TemplateData: map[string]any{
			"Quorum": f.Quorum,
		},
	})
}

func (s *Service) newFormCreateKey(ctx context.Context) (f *FormCreateKey, err error) {
	f = &FormCreateKey{
		Quorum: s.quorumCount,
	}
	f.User, f.lc, err = s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	title, tag, err := f.lc.LocalizeWithTag(&i18n.LocalizeConfig{
		DefaultMessage: createButtonLabel,
	})
	f.Title = title
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

	secretBytes := make([]byte, s.keyLength)
	if _, err = rand.Read(secretBytes); err != nil {
		form.Error = err.Error()
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
		form.Error, err = form.lc.LocalizeMessage(&i18n.Message{
			ID:    "KidwordsErrorTooManyKeys",
			Other: "too many paper keys, please delete one",
		})
		return form, err
	}

	kidwordsSecret, err := kidwords.NewSecret(
		secretBytes,
		s.shardCount,
		s.quorumCount,
	)
	if err != nil {
		return form, err
	}
	fingerPrint := kidwordsSecret.GetFingerprint()
	_, err = rp.Create(ctx, secret.Secret{
		// ID:          id,
		UserID: userID,
		// Name:        base64.RawStdEncoding.EncodeToString(kidwordsSecret.GetFingerprint()),
		Fingerprint: fingerPrint,
		Type:        secret.TypePaperKey,
		SaltedHash:  argonHash.String(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return form, err
	}

	form.Name = fmt.Sprintf("%x", fingerPrint)
	form.Secret = s.encoder.MakeTable(kidwordsSecret)
	return form, nil
}
