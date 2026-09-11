package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FormUpdateKey struct {
	*FormCreateKey
	Name         string
	KeyName      string
	KeyNameLabel string
	KeyNameError string
	keyView
	IsProcessed bool
}

func (f *FormUpdateKey) Title() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateFormTitle",
			Other: "Change Paper Key Name",
		},
	})
}

func (f *FormUpdateKey) Description() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateFormDescription",
			Other: "Change the name of this paper key. Be careful not to betray the exact physical location where it might be stored or which persons might have access to it or knowledge of it.",
		},
	})
}

func (f *FormUpdateKey) NameFieldLabel() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsNameFieldLabel",
			Other: "Name",
		},
	})
}

func (f *FormUpdateKey) UpdateButtonLabel() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateButtonLabel",
			Other: "Update",
		},
	})
}

func (s *Service) updateKeyFormView(ctx context.Context, ID string) (form *FormUpdateKey, err error) {
	form = &FormUpdateKey{}
	form.FormCreateKey, err = s.newFormCreateKey(ctx)
	if err != nil {
		return nil, err
	}
	secret, err := s.repository.Retrieve(ctx, ID)
	if err != nil {
		form.Error = err.Error()
		return
	}
	form.keyView = newKeyView(secret)
	if secret.UserID != form.FormCreateKey.User.GetID() {
		form.Error = newNotFoundError(form.lc).Error()
		return
	}
	return
}

type UpdateKeyRequest struct {
	ID     string
	Name   string
	Method string
}

func (r *UpdateKeyRequest) Validate(context.Context) error {
	if r.Method != http.MethodPost {
		return nil
	}
	if r.ID == "" {
		return errors.New("UUID is required")
	}
	r.Name = strings.TrimSpace(r.Name)
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
	form.ID = req.ID

	if req.Method == http.MethodGet {
		key, err := s.repository.Retrieve(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		if key.UserID != form.User.GetID() {
			return nil, newNotFoundError(form.lc)
		}
		form.Name = key.Name
		return form, nil
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
	form.keyView = newKeyView(key)
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
	form.IsProcessed = true
	return form, nil
}
