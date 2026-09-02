package service

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/staticfs"
)

const (
	TemplateNamePrefix = "kidwords:"
	TemplateHeadName   = TemplateNamePrefix + "head"
	TemplatePageList   = TemplateNamePrefix + "list"
	// TemplatePageName   = TemplateNamePrefix + "page"
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

		head, err := assets.ReadFile("media/head.html")
		if err != nil {
			return fmt.Errorf("unable to load head.html template: %w", err)
		}
		tmpl, err = tmpl.New(TemplateHeadName).Parse(string(head))
		if err != nil {
			return fmt.Errorf("unable to create a <head> template: %w", err)
		}
	}

	return nil
}
