package http

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed media/*
var assets embed.FS

// var TemplateFuncs = map[string]any{
// 	"base64": base64.RawStdEncoding.EncodeToString,
// }

type Templates struct {
	Page   *template.Template
	List   *template.Template
	Create *template.Template
	Update *template.Template
	Delete *template.Template
}

func (t Templates) isComplete() bool {
	return t.Page != nil &&
		t.List != nil &&
		t.Create != nil &&
		t.Update != nil &&
		t.Delete != nil
}

func LoadDefaultTemplatesInto(tmpl *template.Template) (*template.Template, error) {
	data, err := assets.ReadFile("media/templates.html")
	if err != nil {
		return nil, fmt.Errorf("unable to load templates: %w", err)
	}
	tmpl, err = tmpl.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("unable to parse templates: %w", err)
	}
	return tmpl, nil
}

func NewPageTemplate(page, content *template.Template) (*template.Template, error) {
	tmpl, err := page.Clone()
	if err != nil {
		return nil, err
	}
	tmpl, err = tmpl.AddParseTree("content", content.Tree)
	if err != nil {
		return nil, err
	}
	tmpl = tmpl.Lookup(page.Name())
	if tmpl == nil {
		return nil, fmt.Errorf("root page template was not found")
	}
	return tmpl, nil
}
