package http

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/dkotik/htadaptor"
)

type options struct {
	Mux        *http.ServeMux
	PathPrefix string
	Adaptor    *htadaptor.Adaptor
	Templates  *template.Template
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

func WithTemplates(templates *template.Template) Option {
	return func(o *options) error {
		if templates == nil {
			return errors.New("nil templates")
		}
		if o.Templates != nil {
			return errors.New("templates already set")
		}
		o.Templates = templates
		return nil
	}
}
