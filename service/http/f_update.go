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
			Other: "Change Paper Key Name",
		},
	})
}

func (f *formUpdateKey) Description() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsUpdateFormDescription",
			Other: "Change the name of this paper key. Be careful not to betray the exact physical location where it might be stored or which persons might have access to it or knowledge of it.",
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
