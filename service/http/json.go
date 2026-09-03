package http

import (
	"errors"
	"fmt"
	"net/http"

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
		func(o *options) error {
			if o.Mux == nil {
				o.Mux = http.NewServeMux()
			}
			if o.PathPrefix == "" {
				o.PathPrefix = "/"
			}
			if o.Adaptor == nil {
				adaptor := htadaptor.New()
				o.Adaptor = &adaptor
			}
			return nil
		}) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to initialize JSON handler: %w", err)
		}
	}

	create, err := o.Adaptor.AdaptStringFunc(
		s.CreateKeyFormPost,
		extract.StringValueExtractorFunc(func(r *http.Request) (string, error) {
			return r.FormValue("name"), nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create a create key form handler: %w", err)
	}

	update, err := o.Adaptor.AdaptFunc(
		s.UpdateKeyFormPost,
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
