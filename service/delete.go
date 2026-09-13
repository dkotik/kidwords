package service

import (
	"context"
	"fmt"
)

func (s *Service) Delete(ctx context.Context, ID string) (v *KeyView, err error) {
	user, err := s.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate user: %w", err)
	}
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	scrt, err := rp.Retrieve(ctx, ID)
	if err != nil {
		return nil, err
	}
	v, err = NewKeyView(user, scrt)
	if err != nil {
		return nil, err
	}
	return v, rp.Delete(ctx, ID)
}
