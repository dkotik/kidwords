package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dkotik/kidwords/service/secret"
)

const shortDateFormat = "2006-01-02"

type KeyView struct {
	User               User
	ID                 string
	UserID             string
	Name               string
	Type               string
	Fingerprint        string
	CreatedAt          string
	UpdatedAt          string
	LastAcceptedAt     string
	CreatedAtFull      string
	UpdatedAtFull      string
	LastAcceptedAtFull string
}

func NewKeyView(user User, s secret.Secret) (v *KeyView, err error) {
	if s.UserID != user.GetID() {
		return nil, secret.ErrNotFound
	}
	v = &KeyView{
		User:   user,
		ID:     s.ID,
		UserID: s.UserID,
		Name:   s.Name,
		Type:   s.Type,
	}

	v.Fingerprint = base64.RawStdEncoding.EncodeToString(s.Fingerprint)
	v.CreatedAt = s.CreatedAt.Format(shortDateFormat)
	v.CreatedAtFull = s.CreatedAt.Format(time.RFC3339)
	v.UpdatedAt = s.UpdatedAt.Format(shortDateFormat)
	v.UpdatedAtFull = s.UpdatedAt.Format(time.RFC3339)
	if !s.LastAcceptedAt.IsZero() {
		v.LastAcceptedAt = s.LastAcceptedAt.Format(shortDateFormat)
		v.LastAcceptedAtFull = s.LastAcceptedAt.Format(time.RFC3339)
	}
	return v, nil
}

func (s *Service) ViewKey(
	ctx context.Context,
	ID string,
) (*KeyView, error) {
	user, err := s.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate user: %w", err)
	}
	scrt, err := s.repository.Retrieve(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve secret: %w", err)
	}
	return NewKeyView(user, scrt)
}
