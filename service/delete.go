package service

import (
	"context"
	"fmt"
)

func (s *Service) delete(ctx context.Context, ID string) (any, error) {
	user, lc, err := s.unpackContext(ctx)
	if err != nil {
		return nil, err
	}
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	key, err := rp.Retrieve(ctx, ID)
	if err != nil {
		return nil, err
	}
	if key.UserID != user.GetID() {
		return nil, newNotFoundError(lc)
	}
	err = rp.Delete(ctx, ID)
	if err != nil {
		return nil, err
	}
	return s.list(ctx)
}
