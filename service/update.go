package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UpdateKeyRequest struct {
	ID   string
	Name string
}

func (r *UpdateKeyRequest) Validate(context.Context) error {
	if r.ID == "" {
		return errors.New("ID is required")
	}
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.New("Name is required")
	}
	return nil
}

func (s *Service) UpdateKey(ctx context.Context, req *UpdateKeyRequest) (v *KeyView, err error) {
	user, err := s.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate user: %w", err)
	}
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	scrt, err := rp.Retrieve(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	scrt.Name = req.Name
	scrt.UpdatedAt = time.Now()
	err = rp.Update(ctx, scrt)
	if err != nil {
		return nil, err
	}
	return NewKeyView(user, scrt)
}
