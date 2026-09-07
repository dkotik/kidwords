package sqlite

import (
	"context"
	"errors"
	"time"

	"github.com/dkotik/kidwords/service/secret"
	lib "modernc.org/sqlite/lib"
	"zombiezen.com/go/sqlite"
)

func (r *sqRepository) Update(ctx context.Context, s secret.Secret) (err error) {
	defer r.BindContext(ctx)()
	if err = r.stmtUpdate.Reset(); err != nil {
		return err
	}
	if s.Type == "" {
		return errors.New("secret type is required")
	}

	// r.stmtUpdate.BindText(1, s.UserID)
	r.stmtUpdate.BindText(1, s.Name)
	r.stmtUpdate.BindText(2, s.Type)
	r.stmtUpdate.BindText(3, s.SaltedHash)
	r.stmtUpdate.BindText(4, encodeFingerprint(s.Fingerprint))
	r.stmtUpdate.BindText(5, encodeTime(s.CreatedAt))
	r.stmtUpdate.BindText(6, encodeTime(time.Now()))
	r.stmtUpdate.BindText(7, encodeTime(s.LastAcceptedAt))
	r.stmtUpdate.BindText(8, s.ID)

	ok := false
	for {
		ok, err = r.stmtUpdate.Step()
		if err != nil {
			switch code := sqlite.ErrCode(err); code {
			case lib.SQLITE_OK:
				return nil
			case lib.SQLITE_CONSTRAINT_PRIMARYKEY:
				if s.ID == "" {
					return err
				}
				return secret.ErrDuplicateSecretName
			case lib.SQLITE_CONSTRAINT_UNIQUE:
				return secret.ErrDuplicateSecretName
			default:
				return err
			}
		}
		if !ok {
			break
		}
	}
	return nil
}
