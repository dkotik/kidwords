/*
Package service authenticates accounts with paper keys and manages account paper keys.

It is designed to integrate securely with an existing HTTP service.
*/
package service

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"net/http"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/kidwords/service/secret"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed media/*
var assets embed.FS

type User interface {
	GetID() string
	GetName() string
}

type Authenticator interface {
	Authenticate(context.Context) (User, error)
}

type Localizer interface {
	GetLocalizer(context.Context) *i18n.Localizer
}

type LocalizerFunc func(context.Context) *i18n.Localizer

func (f LocalizerFunc) GetLocalizer(ctx context.Context) *i18n.Localizer {
	return f(ctx)
}

type Service struct {
	authenticator Authenticator
	repository    secret.Repository
	localizer     Localizer
	keyCountLimit int
	keyLength     int
	shardCount    int
	quorumCount   int
}

func New(
	authenticator Authenticator,
	repository secret.Repository,
	withOptions ...Option,
) (http.Handler, error) {
	o := &options{}
	for _, opt := range withOptions {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if o.Localizer == nil {
		lc := i18n.NewLocalizer(i18n.NewBundle(language.English))
		o.Localizer = LocalizerFunc(func(context.Context) *i18n.Localizer {
			return lc
		})
	}
	if o.KeyCountLimit == 0 {
		o.KeyCountLimit = DefaultKeyCount
	}
	if o.KeyLength == 0 {
		o.KeyLength = DefaultKeyLength
	}
	if o.ShardCount == 0 {
		o.ShardCount = DefaultShardCount
	}
	if o.QuorumCount == 0 {
		o.QuorumCount = DefaultQuorumCount
	}
	if o.ServeMux == nil {
		o.ServeMux = http.NewServeMux()
	}
	if o.Adaptor == nil {
		o.Adaptor = &htadaptor.Adaptor{}
	}
	if o.Template == nil {
		o.Template = template.Must(template.New("").Parse(""))
	}

	service := &Service{
		authenticator: authenticator,
		repository:    repository,
		localizer:     o.Localizer,
		keyCountLimit: int(o.KeyCountLimit),
		keyLength:     int(o.KeyLength),
		shardCount:    int(o.ShardCount),
		quorumCount:   int(o.QuorumCount),
	}

	list, err := o.Adaptor.AdaptNullaryFunc(
		service.list,
		htadaptor.WithTemplate(nil),
	)
	if err != nil {
		return nil, err
	}
	o.ServeMux.Handle(o.PathPrefix, list)

	return o.ServeMux, nil
}

func (s *Service) unpackContext(ctx context.Context) (u User, l *i18n.Localizer, err error) {
	l = s.localizer.GetLocalizer(ctx)
	if l == nil {
		return nil, nil, errors.New("no localizer in request context")
	}
	u, err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return u, l, nil
}

func newNotFoundError(lc *i18n.Localizer) error {
	msg, err := lc.LocalizeMessage(&i18n.Message{
		ID:    "KidwordsErrorKeyNotFound",
		Other: "paper key does not exist",
	})
	if err != nil {
		return err
	}
	return htadaptor.NewNotFoundError(msg)
}
