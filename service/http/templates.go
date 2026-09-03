package http

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dkotik/htadaptor"
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

var (
	_ htadaptor.Encoder = (*pageRenderer)(nil)
	// buffers                   = &sync.Pool{
	// 	New: func() any {
	// 		return &bytes.Buffer{}
	// 	},
	// }
)

type pageRenderer struct {
	Page *template.Template
	Main *template.Template
}

func NewPageRenderer(page, main *template.Template) htadaptor.Encoder {
	return pageRenderer{
		Page: page,
		Main: main,
	}
}

func (r pageRenderer) Encode(
	w http.ResponseWriter,
	req *http.Request,
	statusCode int,
	v any,
) (err error) {
	tmpl, err := r.Page.Clone()
	if err != nil {
		return err
	}
	tmpl, err = tmpl.AddParseTree("page", r.Page.Tree)
	if err != nil {
		return err
	}
	tmpl, err = tmpl.AddParseTree("content", r.Main.Tree)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; encoding=utf-8")
	w.WriteHeader(statusCode)
	return tmpl.Lookup("page").Execute(w, v)

	// buf := buffers.Get().(*bytes.Buffer)
	// defer buffers.Put(buf)
	// buf.Reset()
	// if err = r.Main.Execute(buf, v); err != nil {
	// 	return err
	// }
	// w.Write(buf.Bytes())
	// return r.Page.Execute(w, struct {
	// 	any
	// 	Content string
	// }{v, buf.String()})
}
