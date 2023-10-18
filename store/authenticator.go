package store

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"log/slog"

	"github.com/google/uuid"
)

// Authenticator simplifies low level [Store] operations and secures them with reasonable defaults. Use it to verify that a particular owner knows one of the [PaperKey]s associated with them.
type Authenticator struct {
	store                    Store
	desiredPaperKeyByteCount int64
}

func NewAuthenticator(using Store, desiredPaperKeyByteCount int64) (*Authenticator, error) {
	if using == nil {
		return nil, errors.New("cannot use a <nil> key store")
	}
	if desiredPaperKeyByteCount < 8 {
		return nil, errors.New("should create keys smaller than 8 bytes")
	}
	if desiredPaperKeyByteCount > 512 {
		return nil, errors.New("should create keys greater than 512 bytes")
	}
	return &Authenticator{
		store:                    using,
		desiredPaperKeyByteCount: desiredPaperKeyByteCount,
	}, nil
}

func (a *Authenticator) Authenticate(
	ctx context.Context,
	keyOwner string,
	key string,
) (bool, error) {
	keys, err := a.store.RetrieveAll(ctx, keyOwner)
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

func (a *Authenticator) AddPaperKey(
	ctx context.Context,
	keyOwner string,
	keyName string,
) (key *PaperKey, err error) {
	secret := &bytes.Buffer{}
	if _, err = io.CopyN(
		secret,
		rand.Reader,
		a.desiredPaperKeyByteCount,
	); err != nil {
		return nil, err
	}

	hash, err := NewArgonHash(secret.Bytes())
	if err != nil {
		return nil, err
	}

	key = &PaperKey{
		ID:         uuid.New().String(),
		Owner:      keyOwner,
		Name:       keyName,
		SaltedHash: hash.String(),
		Created:    time.Now(),
	}
	if err = a.store.Create(ctx, key); err != nil {
		return nil, err
	}
	return key, nil
}

func (a *Authenticator) ListPaperKeys(
	ctx context.Context,
	keyOwner string,
) ([]*PaperKey, error) {
	return a.ListPaperKeys(ctx, keyOwner)
}

func (a *Authenticator) RemovePaperKey(
	ctx context.Context,
	ID string,
	keyOwner string,
) (err error) {
	all, err := a.ListPaperKeys(ctx, keyOwner)
	if err != nil {
		return err
	}
	for _, pk := range all {
		if pk.ID == ID {
			return a.store.Delete(ctx, ID)
		}
	}
	return fmt.Errorf("key does not exist")
}

func (a *Authenticator) RemoveOwner(
	ctx context.Context,
	keyOwner string,
) (err error) {
	return a.store.DeleteByOwner(ctx, keyOwner)
}
