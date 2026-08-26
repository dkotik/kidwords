package htservice

import (
	"os"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestListTemplate(t *testing.T) {
	lc := i18n.NewLocalizer(i18n.NewBundle(language.English), "en")
	if err := templates.Lookup("list").Execute(os.Stdout, &listResponse{
		Secrets: []ArgonSecretLabel{
			{Name: "testKey1", Description: "Description1"},
		},
		lc: lc,
	}); err != nil {
		t.Error(err)
	}
	// t.Fatal("impl")
}
