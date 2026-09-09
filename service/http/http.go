package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/extract"
	"github.com/dkotik/kidwords/service"
)

func New(s *service.Service, withOptions ...Option) (_ http.Handler, err error) {
	if s == nil {
		return nil, errors.New("nil service")
	}
	o := &options{}
	for _, option := range append(
		withOptions,
		withDefaultMux,
		withDefaultPathPrefix,
		withDefaultAdaptor,
		withDefaultTemplates,
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize HTTP handler: %w", err)
		}
	}
	prefix := o.PathPrefix

	create, err := o.Adaptor.AdaptStringFunc(
		s.CreateKeyFormPost,
		extract.StringValueExtractorFunc(func(r *http.Request) (string, error) {
			return r.FormValue("name"), nil
		}),
		htadaptor.WithTemplate(o.Templates.Create),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a create key form handler: %w", err)
	}
	o.Mux.Handle(prefix+"create", create)

	update, err := o.Adaptor.AdaptFunc(
		s.UpdateKeyFormPost,
		htadaptor.WithTemplate(o.Templates.Update),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create an update key form handler: %w", err)
	}
	o.Mux.Handle(prefix+"update", update)

	deletePageTemplate, err := NewPageTemplate(
		o.Templates.Page,
		o.Templates.Delete,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create delete page template: %w", err)
	}
	delete, err := o.Adaptor.AdaptFunc(
		func(ctx context.Context, r *service.DeleteRequest) (any, error) {
			f, err := s.Delete(ctx, r)
			return deleteFormWithRedirect{
				FormDeleteKey: f,
				Location:      prefix,
			}, err
		},
		htadaptor.WithEncoder(
			htadaptor.NewTemporaryRedirect(
				htadaptor.NewTemplateEncoder(deletePageTemplate),
			),
		),
		htadaptor.WithQueryValues("id"),
		// htadaptor.WithTemplate(o.Templates.Delete),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}
	o.Mux.Handle(prefix+"delete", delete)

	listPageTemplate, err := NewPageTemplate(
		o.Templates.Page,
		o.Templates.List,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create list page template: %w", err)
	}
	list, err := o.Adaptor.AdaptNullaryFunc(
		s.List,
		htadaptor.WithTemplate(listPageTemplate),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a list key form handler: %w", err)
	}
	o.Mux.Handle(prefix, list)

	return o.Mux, nil
}

type deleteFormWithRedirect struct {
	*service.FormDeleteKey
	Location string
}

func (f deleteFormWithRedirect) GetRedirectLocation() string {
	if !f.IsProcessed {
		return ""
	}
	return f.Location
}
