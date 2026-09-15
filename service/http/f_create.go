package http

import (
	"context"
	"errors"
	"strings"

	"github.com/dkotik/kidwords/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CreateKeyRequest struct {
	Name  string
	Split uint8
}

func (r *CreateKeyRequest) Validate(ctx context.Context) error {
	return nil
}

func (r *CreateKeyRequest) ValidatePost(ctx context.Context) error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Split == 0 {
		return errors.New("split value is required")
	}
	if r.Split > 12 {
		return errors.New("split value is greater than 12")
	}
	return nil
}

type formCreateKey struct {
	*form
	Name   string
	Split  uint8
	Secret *service.Secret
}

func (f *formCreateKey) Title() (string, error) {
	return f.LabelCreate()
}

func (f *formCreateKey) Description() (string, error) {
	if f.Secret == nil {
		return f.lc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "KidwordsCreateFormPrepareDescription",
				Other: "PREPARE",
			},
		})
	}

	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsCreateFormDescription",
			Other: "Print and store this key in a safe place. The key is split into validatable shards. You will need at least {{ .Quorum }} valid shards to gain access to your account.",
		},
		TemplateData: map[string]any{
			"Quorum": f.Secret.Quorum,
		},
	})
}

func (f *formCreateKey) LabelCreate() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelCreate",
			Other: "Create New Paper Key",
		},
	})
}
