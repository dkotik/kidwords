package store

import (
	"context"
	"errors"

	"log/slog"
)

type Authenticator struct {
	store Store
}

func NewAuthenticator(using Store) (*Authenticator, error) {
	if using == nil {
		return nil, errors.New("cannot use a <nil> key store")
	}
	return &Authenticator{store: using}, nil
}

func (a *Authenticator) Authenticate(
	ctx context.Context,
	owner string,
	key string,
) (bool, error) {
	keys, err := a.store.RetrieveAll(ctx, owner)
	if err != nil {
		return false, err
	}
	for _, key := range keys {
		hash, err := ParseArgonHash(key.SaltedHash)
		if err != nil {
			slog.WarnContext(
				ctx,
				"stored KidWords paper key is corrupted",
				slog.String("ID", key.ID),
				slog.Any("error", err),
			)
			continue
		}
		ok, err := hash.Match([]byte(key.SaltedHash))
		if err != nil {
			slog.WarnContext(
				ctx,
				"stored KidWords paper key matching failed",
				slog.String("ID", key.ID),
				slog.Any("error", err),
			)
			continue
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
