package http

import (
	"github.com/dkotik/kidwords/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type formCreateKey struct {
	*form
	*service.Secret
}

func (f *formCreateKey) Title() (string, error) {
	return f.LabelCreate()
}

func (f *formCreateKey) Description() (string, error) {
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

func (f *formCreateKey) LabelCreate() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelCreate",
			Other: "Create New Paper Key",
		},
	})
}
