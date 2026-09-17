package http

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type form struct {
	lc          *i18n.Localizer
	Locale      string
	LabelCancel string
	Redirect    string
	Error       string
}

func newForm(lc *i18n.Localizer) (f *form, err error) {
	f = &form{
		lc: lc,
	}
	var tag language.Tag
	f.LabelCancel, tag, err = lc.LocalizeWithTag(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelCancel",
			Other: "Cancel",
		},
	})
	if err != nil {
		return
	}
	f.Locale = tag.String()
	return f, nil
}

func (f *form) GetRedirectLocation() string {
	return f.Redirect
}

func (f *form) LabelSaveChanges() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelSaveChanges",
			Other: "Save Changes",
		},
	})
}

func (f *form) LabelBack() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelBack",
			Other: "Back",
		},
	})
}

func (f *form) LabelClose() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelClose",
			Other: "Close",
		},
	})
}

func (f *form) LabelDelete() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelDelete",
			Other: "Delete",
		},
	})
}

func (f *form) NameFieldLabel() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsNameFieldLabel",
			Other: "Paper Key Name",
		},
	})
}

func (f *form) LabelAddCustomName() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelAddCustomName",
			Other: "Add a Custom Name",
		},
	})
}

func (f *form) LabelAddCustomNameHelp() (string, error) {
	return f.lc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "KidwordsLabelAddCustomNameHelp",
			Other: "Be careful not to betray the exact physical location where the paper key might be stored or which persons might have access to it or knowledge of it.",
		},
	})
}
