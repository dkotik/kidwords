package http

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/extract"
	"github.com/dkotik/htadaptor/staticfs"
	"github.com/dkotik/kidwords/service"
)

const (
	TemplateNamePrefix = "kidwords/"
	TemplateHeadName   = TemplateNamePrefix + "header.html"
	TemplatePageList   = TemplateNamePrefix + "list.html"
	TemplateCreatePage = TemplateNamePrefix + "create.html"
	TemplateUpdatePage = TemplateNamePrefix + "update.html"
	TemplateListPage   = TemplateNamePrefix + "list.html"
)

//go:embed media/*
var assets embed.FS

func LoadDefaultTemplates() (*template.Template, error) {
	data, err := assets.ReadFile("media/templates.html")
	if err != nil {
		return nil, fmt.Errorf("unable to load templates: %w", err)
	}
	tmpl, err := template.New("").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("unable to parse templates: %w", err)
	}
	return tmpl, nil
}

func New(s *service.Service, withOptions ...Option) (_ http.Handler, err error) {
	if s == nil {
		return nil, errors.New("nil service")
	}
	o := &options{}
	for _, option := range append(
		withOptions,
		func(o *options) error {
			if o.Mux == nil {
				o.Mux = http.NewServeMux()
			}
			if o.PathPrefix == "" {
				o.PathPrefix = "/"
			}
			if o.Templates == nil {
				o.Templates, err = LoadDefaultTemplates()
				if err != nil {
					return err
				}
			}
			for _, name := range []string{
				TemplatePageList,
				TemplateCreatePage,
				TemplateUpdatePage,
				TemplateListPage,
			} {
				if tmpl := o.Templates.Lookup(name); tmpl == nil {
					return fmt.Errorf("template %s not found", name)
				}
			}
			head := o.Templates.Lookup(TemplateHeadName)
			if head == nil {
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

				head, err := assets.ReadFile("media/header.html")
				if err != nil {
					return fmt.Errorf("unable to load header.html template: %w", err)
				}
				o.Templates, err = o.Templates.New(TemplateHeadName).Parse(string(head))
				if err != nil {
					return fmt.Errorf("unable to create a <head> template: %w", err)
				}
			}

			if o.Adaptor == nil {
				adaptor := htadaptor.New()
				o.Adaptor = &adaptor
			}
			return nil
		}) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize HTTP handler: %w", err)
		}
	}

	create, err := o.Adaptor.AdaptStringFunc(
		s.CreateKeyFormPost,
		extract.StringValueExtractorFunc(func(r *http.Request) (string, error) {
			return r.FormValue("name"), nil
		}),
		htadaptor.WithTemplate(o.Templates.Lookup(TemplateCreatePage)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a create key form handler: %w", err)
	}

	update, err := o.Adaptor.AdaptFunc(
		s.UpdateKeyFormPost,
		htadaptor.WithTemplate(o.Templates.Lookup(TemplateUpdatePage)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create an update key form handler: %w", err)
	}

	idExtractor, err := extract.NewQueryValueExtractor("delete")
	if err != nil {
		return nil, err
	}
	delete, err := o.Adaptor.AdaptStringFunc(
		s.Delete,
		idExtractor,
		htadaptor.WithTemplate(o.Templates.Lookup(TemplateListPage)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}

	list, err := o.Adaptor.AdaptNullaryFunc(
		s.List,
		htadaptor.WithTemplate(o.Templates.Lookup(TemplateListPage)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a list key form handler: %w", err)
	}

	o.Mux.Handle(o.PathPrefix, htadaptor.NewMethodMux(&htadaptor.MethodSwitch{
		Get:    list,
		Post:   create,
		Put:    update,
		Delete: delete,
	}))

	return o.Mux, nil
}
