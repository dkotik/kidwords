package sqlite

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/dkotik/kidwords/service/secret"
	lib "modernc.org/sqlite/lib"
	"zombiezen.com/go/sqlite"
)

func (r *sqRepository) Create(ctx context.Context, s secret.Secret) (ID string, err error) {
	defer r.BindContext(ctx)()
	if err = r.stmtCreate.Reset(); err != nil {
		return ID, err
	}
	if s.ID == "" {
		ID = uuid.New().String()
		s.ID = ID
	} else {
		ID = s.ID
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = s.CreatedAt
	}
	if s.Type == "" {
		return ID, errors.New("secret type is required")
	}

	r.stmtCreate.BindText(1, s.ID)
	r.stmtCreate.BindText(2, s.UserID)
	r.stmtCreate.BindText(3, s.Name)
	r.stmtCreate.BindText(4, secret.TypePaperKey)
	r.stmtCreate.BindText(5, s.SaltedHash)
	r.stmtCreate.BindText(6, encodeTime(s.CreatedAt))
	r.stmtCreate.BindText(7, encodeTime(s.UpdatedAt))

	ok := false
	for {
		ok, err = r.stmtCreate.Step()
		if err != nil {
			switch code := sqlite.ErrCode(err); code {
			case lib.SQLITE_OK:
				return ID, nil
			case lib.SQLITE_CONSTRAINT_PRIMARYKEY:
				if s.ID == "" {
					return ID, err
				}
				return s.ID, secret.ErrDuplicateSecretName
			case lib.SQLITE_CONSTRAINT_UNIQUE:
				return s.ID, secret.ErrDuplicateSecretName
			default:
				return ID, err
			}
		}
		if !ok {
			break
		}
	}
	return ID, nil
}
