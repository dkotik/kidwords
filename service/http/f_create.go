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
	// if r.Name == "" {
	// 	return errors.New("name is required")
	// }
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
	return f.form.LabelCreate()
}

func (f *formCreateKey) Description() (string, error) {
	if f.Secret == nil {
		return f.lc.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "KidwordsCreateFormPrepareDescription",
				Other: "Print and store the new paper key in a safe place. You may use it in the future to gain access to your account in situations when all other methods of authentication are unavailable. System administrators may retain a copy of the new key for a short time as proof of account ownership even when deleted. If a hacker gains access to your account and replaces all existing paper keys with new ones, you may be able to regain access using an old key.",
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

func (f *formCreateKey) LabelCopy() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelCopy",
			Other: "Copy",
		},
	})
}

func (f *formCreateKey) LabelPrint() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelPrint",
			Other: "Print",
		},
	})
}

func (f *formCreateKey) LabelSplit() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelSplit",
			Other: "Split Into Multiple Printable Pages",
		},
	})
}

func (f *formCreateKey) LabelCloseWarning() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelCloseWarning",
			Other: "Are you ready to close the paper key? If you did not print or save it, you will never be able to use it.",
		},
	})
}

func (f *formCreateKey) LabelSplitHelp() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelSplitHelp",
			Other: "If you split the key into multiple pages, store each page in a separate secure location or with a trusted third party.",
		},
	})
}
