package sqlite

import (
	"context"

	"github.com/dkotik/kidwords/service/secret"
)

func (s *sqRepository) Create(ctx context.Context, p secret.Secret) error {
	// _, err := s.stmtCreate.ExecContext(
	// 	ctx,
	// 	p.ID,
	// 	p.Owner,
	// 	p.Name,
	// 	p.SaltedHash,
	// 	p.Created.Unix(),
	// )
	return nil
}
