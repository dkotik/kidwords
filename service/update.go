package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FormUpdateKey struct {
	lc           *i18n.Localizer
	User         User
	UUID         string
	KeyName      string
	KeyNameError string
}

func (s *Service) updateKeyFormView(ctx context.Context, UUID string) (any, error) {
	user, lc, err := s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	secret, err := s.repository.Retrieve(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if secret.UserID != user.GetID() {
		return nil, newNotFoundError(lc)
	}
	return &FormUpdateKey{
		lc:      lc,
		User:    user,
		UUID:    UUID,
		KeyName: secret.Name,
	}, nil
}

type UpdateKeyRequest struct {
	UUID string
	Name string
}

func (r *UpdateKeyRequest) Validate() error {
	if r.UUID == "" {
		return errors.New("UUID is required")
	}
	if r.Name == "" {
		return errors.New("Name is required")
	}
	return nil
}

func (s *Service) updateKeyFormPost(ctx context.Context, req *UpdateKeyRequest) (_ *FormUpdateKey, err error) {
	form := &FormUpdateKey{}
	user, lc, err := s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	current, err := s.repository.Retrieve(ctx, req.UUID)
	if err != nil {
		return nil, err
	}
	if current.UserID != user.GetID() {
		return nil, newNotFoundError(lc)
	}
	form.KeyName = strings.TrimSpace(req.Name)
	if form.KeyName == "" {
		form.KeyNameError, err = form.lc.LocalizeMessage(&i18n.Message{
			ID:    "KidwordsErrorKeyNameRequired",
			Other: "key name is empty",
		})
		return form, err
	}

	err = s.repository.Update(ctx, secret.Secret{
		ID:         req.UUID,
		Name:       form.KeyName,
		Type:       secret.TypePaperKey,
		SaltedHash: current.SaltedHash,
		CreatedAt:  current.CreatedAt,
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return form, err
	}

	return form, nil
}
