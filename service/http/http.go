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
		withDefaultLocalizer,
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize HTTP handler: %w", err)
		}
	}
	prefix := o.PathPrefix
	lcExtractor := o.Localizer
	idExtractor, err := extract.NewQueryValueExtractor("id")
	if err != nil {
		return nil, fmt.Errorf("unable to create query value extractor for delete key form: %w", err)
	}

	createPageTemplate, err := NewPageTemplate(
		o.Templates.Page,
		o.Templates.Create,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create page creation template: %w", err)
	}
	create, err := o.Adaptor.AdaptStringFunc(
		func(ctx context.Context, name string) (*formCreateKey, error) {
			scrt, err := s.CreateKey(ctx, name)
			if err != nil {
				return nil, err
			}
			lc := lcExtractor.GetLocalizer(ctx)
			form, err := newForm(lc)
			if err != nil {
				return nil, fmt.Errorf("unable to create form: %w", err)
			}
			return &formCreateKey{
				form:   form,
				Secret: scrt,
			}, nil
		},
		extract.StringValueExtractorFunc(func(r *http.Request) (string, error) {
			return r.FormValue("name"), nil
		}),
		htadaptor.WithTemplate(createPageTemplate),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a create key form handler: %w", err)
	}
	o.Mux.Handle(prefix+"create", create)

	updatePageTemplate, err := NewPageTemplate(
		o.Templates.Page,
		o.Templates.Update,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create update page template: %w", err)
	}

	getUpdate, err := o.Adaptor.AdaptStringFunc(
		func(ctx context.Context, ID string) (*formUpdateKey, error) {
			keyView, err := s.ViewKey(ctx, ID)
			if err != nil {
				return nil, err
			}
			lc := lcExtractor.GetLocalizer(ctx)
			form, err := newForm(lc)
			if err != nil {
				return nil, fmt.Errorf("unable to create form: %w", err)
			}
			return &formUpdateKey{
				form:    form,
				KeyView: keyView,
			}, nil
		},
		idExtractor,
		htadaptor.WithTemplate(updatePageTemplate),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}
	o.Mux.Handle(prefix+"update", getUpdate)

	update, err := o.Adaptor.AdaptFunc(
		func(ctx context.Context, r *service.UpdateKeyRequest) (*formUpdateKey, error) {
			keyView, err := s.UpdateKey(ctx, r)
			if err != nil {
				return nil, err
			}
			lc := lcExtractor.GetLocalizer(ctx)
			form, err := newForm(lc)
			if err != nil {
				return nil, fmt.Errorf("unable to create form: %w", err)
			}
			form.Redirect = prefix
			return &formUpdateKey{
				form:    form,
				KeyView: keyView,
			}, nil
		},
		htadaptor.WithEncoder(
			htadaptor.NewTemporaryRedirect(
				htadaptor.NewTemplateEncoder(updatePageTemplate),
			),
		),
		// htadaptor.WithQueryValues("id"),
		// htadaptor.WithTemplate(o.Templates.Delete),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}
	o.Mux.Handle("POST "+prefix+"update", update)

	deletePageTemplate, err := NewPageTemplate(
		o.Templates.Page,
		o.Templates.Delete,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create delete page template: %w", err)
	}

	getDelete, err := o.Adaptor.AdaptStringFunc(
		func(ctx context.Context, ID string) (*formDeleteKey, error) {
			keyView, err := s.ViewKey(ctx, ID)
			if err != nil {
				return nil, err
			}
			lc := lcExtractor.GetLocalizer(ctx)
			form, err := newForm(lc)
			if err != nil {
				return nil, fmt.Errorf("unable to create form: %w", err)
			}
			return &formDeleteKey{
				form:    form,
				KeyView: keyView,
			}, nil
		},
		idExtractor,
		htadaptor.WithTemplate(deletePageTemplate),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}
	o.Mux.Handle(prefix+"delete", getDelete)

	delete, err := o.Adaptor.AdaptStringFunc(
		func(ctx context.Context, ID string) (*formDeleteKey, error) {
			keyView, err := s.Delete(ctx, ID)
			if err != nil {
				return nil, err
			}
			lc := lcExtractor.GetLocalizer(ctx)
			form, err := newForm(lc)
			if err != nil {
				return nil, fmt.Errorf("unable to create form: %w", err)
			}
			form.Redirect = prefix
			return &formDeleteKey{
				form:    form,
				KeyView: keyView,
			}, nil
		},
		extract.StringValueExtractorFunc(
			func(r *http.Request) (string, error) {
				return r.FormValue("id"), nil
			},
		),
		htadaptor.WithEncoder(
			htadaptor.NewTemporaryRedirect(
				htadaptor.NewTemplateEncoder(deletePageTemplate),
			),
		),
		// htadaptor.WithQueryValues("id"),
		// htadaptor.WithTemplate(o.Templates.Delete),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}
	o.Mux.Handle("POST "+prefix+"delete", delete)

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
