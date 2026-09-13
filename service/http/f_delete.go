package http

import (
	"github.com/dkotik/kidwords/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type formDeleteKey struct {
	*form
	*service.KeyView
}

func (f *formDeleteKey) Title() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsDeleteFormTitle",
			Other: "Delete Paper Key",
		},
	})
}

func (f *formDeleteKey) Description() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsDeleteFormDescription",
			Other: "Are you certain that you would like to remove this paper key from your account? You will not be able to restore access to your account using this key anymore. Service administration might keep a copy of the key for some time as proof of account ownership, by service policy, in case your account has been hacked and access to it must be restored.",
		},
	})
}
