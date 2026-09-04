package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/staticfs"
)

type options struct {
	Mux        *http.ServeMux
	PathPrefix string
	Adaptor    *htadaptor.Adaptor
	Templates  *Templates
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

func withDefaultMux(o *options) error {
	if o.Mux == nil {
		o.Mux = http.NewServeMux()
	}
	return nil
}

func withDefaultPathPrefix(o *options) error {
	if o.PathPrefix == "" {
		o.PathPrefix = "/"
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
	templates, err := LoadDefaultTemplates()
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
