package http

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/staticfs"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer interface {
	GetLocalizer(context.Context) *i18n.Localizer
}

type LocalizerFunc func(context.Context) *i18n.Localizer

func (f LocalizerFunc) GetLocalizer(ctx context.Context) *i18n.Localizer {
	return f(ctx)
}

type options struct {
	Mux        *http.ServeMux
	PathPrefix string
	Adaptor    *htadaptor.Adaptor
	Templates  *Templates
	Localizer  Localizer
}

type Option func(*options) error

func WithServeMux(mux *http.ServeMux, pathPrefix string) Option {
	return func(o *options) (err error) {
		if mux == nil {
			return errors.New("nil http.ServeMux")
		}
		if err = WithPathPrefix(pathPrefix)(o); err != nil {
			return err
		}
		if o.Mux != nil {
			return errors.New("mux already set")
		}
		o.Mux = mux
		return nil
	}
}

func WithPathPrefix(prefix string) Option {
	return func(o *options) error {
		if prefix == "" {
			return errors.New("empty path prefix")
		}
		if o.PathPrefix != "" {
			return errors.New("path prefix already set")
		}
		o.PathPrefix = prefix
		return nil
	}
}

func WithAdaptor(adaptor htadaptor.Adaptor) Option {
	return func(o *options) error {
		if o.Adaptor != nil {
			return errors.New("adaptor already set")
		}
		o.Adaptor = &adaptor
		return nil
	}
}

func WithTemplates(templates Templates) Option {
	return func(o *options) error {
		if o.Templates != nil {
			return errors.New("templates already set")
		}
		o.Templates = &templates
		return nil
	}
}

func WithLocalizer(localizer Localizer) Option {
	return func(o *options) error {
		if localizer == nil {
			return errors.New("nil localizer")
		}
		if o.Localizer != nil {
			return errors.New("localizer already set")
		}
		o.Localizer = localizer
		return nil
	}
}

func withDefaultMux(o *options) error {
	if o.Mux == nil {
		o.Mux = http.NewServeMux()
	}
	return nil
}

func withDefaultPathPrefix(o *options) error {
	if !strings.HasSuffix(o.PathPrefix, "/") {
		o.PathPrefix = o.PathPrefix + "/"
	}
	return nil
}

func withDefaultAdaptor(o *options) error {
	if o.Adaptor == nil {
		o.Adaptor = &htadaptor.Adaptor{}
	}
	return nil
}

func withDefaultTemplates(o *options) error {
	if o.Templates == nil {
		o.Templates = &Templates{}
	} else if o.Templates.isComplete() {
		return nil // all good
	}
	templates, err := LoadDefaultTemplatesInto(template.New(""))
	if err != nil {
		return err
	}

	if o.Templates.Page == nil {
		o.Templates.Page = templates.Lookup("page")
		htmx, err := assets.ReadFile("media/htmx.min.js")
		if err != nil {
			return fmt.Errorf("unable to load HTMX source: %w", err)
		}
		htmxPath := o.PathPrefix + "htmx.min.js"
		o.Mux.Handle(htmxPath, staticfs.NewFastFileSystemFileWithContentType(htmx, "text/javascript"))

		bulma, err := assets.ReadFile("media/bulma.min.css")
		if err != nil {
			return fmt.Errorf("unable to load Bulma source: %w", err)
		}
		bulmaPath := o.PathPrefix + "bulma.min.css"
		o.Mux.Handle(bulmaPath, staticfs.NewFastFileSystemFileWithContentType(bulma, "text/css"))
	}
	if o.Templates.Create == nil {
		o.Templates.Create = templates.Lookup("create")
	}
	if o.Templates.Update == nil {
		o.Templates.Update = templates.Lookup("update")
	}
	if o.Templates.Delete == nil {
		o.Templates.Delete = templates.Lookup("delete")
	}
	if o.Templates.List == nil {
		o.Templates.List = templates.Lookup("list")
	}

	for _, name := range []string{
		"page",
		"create",
		"update",
		"delete",
		"list",
	} {
		if tmpl := templates.Lookup(name); tmpl == nil {
			return fmt.Errorf("template %s not found", name)
		}
	}
	if !o.Templates.isComplete() {
		return errors.New("some default templates are missing")
	}
	return nil
}

func withDefaultLocalizer(o *options) error {
	if o.Localizer == nil {
		lc := i18n.NewLocalizer(i18n.NewBundle(language.English))
		o.Localizer = LocalizerFunc(func(context.Context) *i18n.Localizer {
			return lc
		})
	}
	return nil
}
