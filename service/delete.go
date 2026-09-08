package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var deleteButtonLabel = &i18n.Message{
	ID:    "KidwordsDeleteButtonLabel",
	Other: "Delete",
}

type FormDeleteKey struct {
	lc   *i18n.Localizer
	User User
	keyView
	Error string
}

func (f *FormDeleteKey) Title() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsDeleteFormTitle",
			Other: "Delete Paper Key",
		},
	})
}

func (f *FormDeleteKey) Description() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsDeleteFormDescription",
			Other: "Are you certain that you would like to remove this paper key from your account? You will not be able to restore access to your account using this key anymore. Service administration might keep a copy of the key for some time as proof of account ownership, by service policy, in case your account has been hacked and access to it must be restored.",
		},
	})
}

func (f *FormDeleteKey) DeleteButtonLabel() (string, error) {
	return f.lc.LocalizeMessage(deleteButtonLabel)
}

type DeleteRequest struct {
	ID           string
	Confirmation bool
}

func (r *DeleteRequest) Validate(ctx context.Context) error {
	if r.ID == "" {
		return errors.New("empty ID")
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, r *DeleteRequest) (form *FormDeleteKey, err error) {
	form = &FormDeleteKey{}
	form.User, form.lc, err = s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	key, err := rp.Retrieve(ctx, r.ID)
	if err != nil {
		form.Error = err.Error()
		return
	}
	form.keyView = newKeyView(key)
	if key.UserID != form.User.GetID() {
		form.Error = newNotFoundError(form.lc).Error()
		return
	}
	err = rp.Delete(ctx, r.ID)
	if err != nil {
		form.Error = err.Error()
		return
	}
	return
}
