package http

import (
	"github.com/dkotik/kidwords/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type formUpdateKey struct {
	*form
	*service.KeyView
}

func (f *formUpdateKey) Title() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateFormTitle",
			Other: "Update Paper Key",
		},
	})
}

func (f *formUpdateKey) Description() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateFormDescription",
			Other: "Update paper key.",
		},
	})
}

func (f *formUpdateKey) UpdateButtonLabel() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateButtonLabel",
			Other: "Update",
		},
	})
}
