package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/htadaptor/extract"
	"github.com/dkotik/kidwords/service"
)

func NewJSON(s *service.Service, withOptions ...Option) (_ http.Handler, err error) {
	if s == nil {
		return nil, errors.New("nil service")
	}
	o := &options{}
	for _, option := range append(
		withOptions,
		withDefaultMux,
		withDefaultPathPrefix,
		withDefaultAdaptor,
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize JSON handler: %w", err)
		}
	}
	idExtractor := extract.StringValueExtractorFunc(
		func(r *http.Request) (string, error) {
			// DELETE requests form body is not automatically parsed
			// by the standard library, so we need to read it manually.
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return "", err
			}
			defer r.Body.Close()

			formData, err := url.ParseQuery(string(body))
			if err != nil {
				return "", err

			}
			return formData.Get("id"), nil
		},
	)

	create, err := o.Adaptor.AdaptStringFunc(
		s.CreateKey,
		extract.StringValueExtractorFunc(func(r *http.Request) (string, error) {
			return r.FormValue("name"), nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a create key form handler: %w", err)
	}

	update, err := o.Adaptor.AdaptFunc(
		s.UpdateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create an update key form handler: %w", err)
	}

	delete, err := o.Adaptor.AdaptStringFunc(
		s.Delete,
		idExtractor,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a delete key form handler: %w", err)
	}

	list, err := o.Adaptor.AdaptNullaryFunc(
		s.List,
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
