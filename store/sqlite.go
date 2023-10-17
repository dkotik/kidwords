package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const InstallSQLiteTable = `
    CREATE TABLE IF NOT EXISTS %s (
      id         TEXT PRIMARY KEY,
      owner      TEXT NOT NULL,
      name       TEXT NOT NULL,
      saltedHash TEXT NOT NULL,
      created    INTEGER NOT NULL
    );`

var _ Store = (*SQLiteStore)(nil) // interface satisfaction

type SQLiteStore struct {
	stmtCreate        *sql.Stmt
	stmtRetrieveAll   *sql.Stmt
	stmtDelete        *sql.Stmt
	stmtDeleteByOwner *sql.Stmt
}

func NewSQLiteStore(database *sql.DB, withOptions ...Option) (s *SQLiteStore, err error) {
	if database == nil {
		return nil, errors.New("cannot use a <nil> database connection")
	}
	o := &options{}
	for _, option := range append(withOptions, withDefaultTableName()) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize the store: %w", err)
		}
	}

	if o.CreateTable {
		if _, err = database.Exec(fmt.Sprintf(InstallSQLiteTable, o.TableName)); err != nil {
			return nil, fmt.Errorf("cannot create the store table: %w", err)
		}
	}

	s = &SQLiteStore{}

	if s.stmtCreate, err = database.Prepare(
		fmt.Sprintf(`
      INSERT INTO %s(id, owner, name, saltedHash, created) VALUES($1, $2, $3, $4, $5);
    `, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtRetrieveAll, err = database.Prepare(
		fmt.Sprintf(`
      SELECT
        id, owner, name, saltedHash, created
      FROM %s WHERE owner=$1 ORDER BY created DESC;`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtDelete, err = database.Prepare(
		fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtDeleteByOwner, err = database.Prepare(
		fmt.Sprintf(`DELETE FROM %s WHERE owner=$1;`, o.TableName),
	); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLiteStore) Create(ctx context.Context, p *PaperKey) error {
	_, err := s.stmtCreate.ExecContext(
		ctx,
		p.ID,
		p.Owner,
		p.Name,
		p.SaltedHash,
		p.Created.Unix(),
	)
	return err
}

func (s *SQLiteStore) RetrieveAll(
	ctx context.Context,
	owner string,
) (result []*PaperKey, err error) {
	rows, err := s.stmtRetrieveAll.QueryContext(ctx, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var created int64
	for rows.Next() {
		key := &PaperKey{}
		if err := rows.Scan(
			&key.ID,
			&key.Owner,
			&key.Name,
			&key.SaltedHash,
			&created,
		); err != nil {
			return nil, err
		}
		key.Created = time.Unix(created, 0)
		result = append(result, key)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SQLiteStore) Delete(ctx context.Context, ID string) error {
	_, err := s.stmtDelete.ExecContext(ctx, ID)
	return err
}

func (s *SQLiteStore) DeleteByOwner(ctx context.Context, owner string) error {
	_, err := s.stmtDeleteByOwner.ExecContext(ctx, owner)
	return err
}
