package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FormUpdateKey struct {
	*FormCreateKey
	UUID string
}

func (f *FormUpdateKey) Title() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateFormTitle",
			Other: "Update Paper Key",
		},
	})
}

func (s *Service) updateKeyFormView(ctx context.Context, UUID string) (any, error) {
	secret, err := s.repository.Retrieve(ctx, UUID)
	if err != nil {
		return nil, err
	}
	f, err := s.newFormCreateKey(ctx)
	if err != nil {
		return nil, err
	}
	if secret.UserID != f.User.GetID() {
		return nil, newNotFoundError(f.lc)
	}
	return &FormUpdateKey{
		FormCreateKey: f,
		UUID:          UUID,
	}, nil
}

type UpdateKeyRequest struct {
	ID   string
	Name string
}

func (r *UpdateKeyRequest) Validate(context.Context) error {
	if r.ID == "" {
		return errors.New("UUID is required")
	}
	if r.Name == "" {
		return errors.New("Name is required")
	}
	return nil
}

func (s *Service) UpdateKeyFormPost(ctx context.Context, req *UpdateKeyRequest) (form *FormUpdateKey, err error) {
	form = &FormUpdateKey{}
	form.FormCreateKey, err = s.newFormCreateKey(ctx)
	if err != nil {
		return nil, err
	}

	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return form, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	key, err := rp.Retrieve(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if key.UserID != form.User.GetID() {
		return nil, newNotFoundError(form.lc)
	}
	form.KeyName = strings.TrimSpace(req.Name)
	if form.KeyName == "" {
		form.KeyNameError, err = form.lc.LocalizeMessage(&i18n.Message{
			ID:    "KidwordsErrorKeyNameRequired",
			Other: "key name is empty",
		})
		return form, err
	}

	key.Name = form.KeyName
	key.UpdatedAt = time.Now()
	err = rp.Update(ctx, key)
	if err != nil {
		return form, err
	}

	return form, nil
}
