/*
Package kwsql persists paper keys to an SQL database.
*/
package kwsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dkotik/kidwords/store"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	stmtCreate      *sql.Stmt
	stmtRetrieveAll *sql.Stmt
	stmtDelete      *sql.Stmt
}

func NewStore(database *sql.DB, tableName string) (s *Store, err error) {
	s = &Store{}

	if s.stmtCreate, err = database.Prepare(
		fmt.Sprintf(`
      INSERT INTO %s(id, owner, name, saltedHash, created) VALUES($1, $2, $3, $4, $5);
    `, tableName),
	); err != nil {
		return nil, err
	}

	if s.stmtRetrieveAll, err = database.Prepare(
		fmt.Sprintf(`
      SELECT
        id, owner, name, saltedHash, created
      FROM %s WHERE owner=$1;`, tableName),
	); err != nil {
		return nil, err
	}

	if s.stmtDelete, err = database.Prepare(
		fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, tableName),
	); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Create(ctx context.Context, p *store.PaperKey) error {
	_, err := s.stmtCreate.ExecContext(
		ctx,
		p.ID,
		p.Owner,
		p.Name,
		p.SaltedHash,
		p.Created,
	)
	return err
}

func (s *Store) RetrieveAll(
	ctx context.Context,
	owner string,
) (result []*store.PaperKey, err error) {
	rows, err := s.stmtCreate.QueryContext(ctx, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		key := &store.PaperKey{}
		if err := rows.Scan(
			&key.ID,
			&key.Owner,
			&key.Name,
			&key.SaltedHash,
			&key.Created,
		); err != nil {
			return nil, err
		}
		result = append(result, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Store) Delete(ctx context.Context, ID string) error {
	_, err := s.stmtDelete.ExecContext(ctx, ID)
	return err
}
