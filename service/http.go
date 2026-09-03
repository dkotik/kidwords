package service

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/extract"
	"github.com/dkotik/htadaptor/staticfs"
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

func (s *Service) MountMux(
	mux *http.ServeMux,
	adaptor htadaptor.Adaptor,
	pathPrefix string,
	tmpl *template.Template,
) (err error) {
	if mux == nil {
		return errors.New("nil HTTP multiplexer")
	}
	if pathPrefix == "" {
		pathPrefix = "/"
	}
	if tmpl == nil {
		tmpl, err = LoadDefaultTemplates()
		if err != nil {
			return err
		}
	}

	head := tmpl.Lookup(TemplateHeadName)
	if head == nil {
		htmx, err := assets.ReadFile("media/htmx.min.js")
		if err != nil {
			return fmt.Errorf("unable to load HTMX source: %w", err)
		}
		htmxPath := pathPrefix + "htmx.min.js"
		mux.Handle(htmxPath, staticfs.NewFastFileSystemFileWithContentType(htmx, "text/javascript"))

		bulma, err := assets.ReadFile("media/bulma.min.css")
		if err != nil {
			return fmt.Errorf("unable to load Bulma source: %w", err)
		}
		bulmaPath := pathPrefix + "bulma.min.css"
		mux.Handle(bulmaPath, staticfs.NewFastFileSystemFileWithContentType(bulma, "text/css"))

		head, err := assets.ReadFile("media/header.html")
		if err != nil {
			return fmt.Errorf("unable to load header.html template: %w", err)
		}
		tmpl, err = tmpl.New(TemplateHeadName).Parse(string(head))
		if err != nil {
			return fmt.Errorf("unable to create a <head> template: %w", err)
		}
	}

	create, err := adaptor.AdaptStringFunc(
		s.createKeyFormPost,
		extract.StringValueExtractorFunc(func(r *http.Request) (string, error) {
			return r.FormValue("name"), nil
		}),
		htadaptor.WithTemplate(tmpl.Lookup(TemplateCreatePage)),
	)
	if err != nil {
		return fmt.Errorf("unable to create a create key form handler: %w", err)
	}

	update, err := adaptor.AdaptFunc(
		s.updateKeyFormPost,
		htadaptor.WithTemplate(tmpl.Lookup(TemplateUpdatePage)),
	)
	if err != nil {
		return fmt.Errorf("unable to create a create key form handler: %w", err)
	}

	idExtractor, err := extract.NewQueryValueExtractor("delete")
	if err != nil {
		return err
	}
	delete, err := adaptor.AdaptStringFunc(
		s.delete,
		idExtractor,
		htadaptor.WithTemplate(tmpl.Lookup(TemplateListPage)),
	)
	if err != nil {
		return fmt.Errorf("unable to create a delete key form handler: %w", err)
	}

	list, err := adaptor.AdaptNullaryFunc(
		s.list,
		htadaptor.WithTemplate(tmpl.Lookup(TemplateListPage)),
	)
	if err != nil {
		return fmt.Errorf("unable to create a create key form handler: %w", err)
	}

	mux.Handle(pathPrefix, htadaptor.NewMethodMux(&htadaptor.MethodSwitch{
		Get:    list,
		Post:   create,
		Put:    update,
		Delete: delete,
	}))

	return nil
}
