package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const InstallPostgresTable = `
    CREATE TABLE IF NOT EXISTS %s (
      id         UUID PRIMARY KEY,
      owner      UUID NOT NULL,
      name       VARCHAR(128) NOT NULL,
      saltedHash TEXT NOT NULL,
      created    TIMESTAMP NOT NULL DEFAULT NOW()
    );`

var _ Store = (*PostgresStore)(nil) // interface satisfaction

type PostgresStore struct {
	stmtCreate        *sql.Stmt
	stmtRetrieveAll   *sql.Stmt
	stmtDelete        *sql.Stmt
	stmtDeleteByOwner *sql.Stmt
}

func NewPostgresStore(database *sql.DB, withOptions ...Option) (s *PostgresStore, err error) {
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
		if _, err = database.Exec(fmt.Sprintf(InstallPostgresTable, o.TableName)); err != nil {
			return nil, fmt.Errorf("cannot create the store table: %w", err)
		}
	}

	s = &PostgresStore{}

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

func (s *PostgresStore) Create(ctx context.Context, p *PaperKey) error {
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

func (s *PostgresStore) RetrieveAll(
	ctx context.Context,
	owner string,
) (result []*PaperKey, err error) {
	rows, err := s.stmtRetrieveAll.QueryContext(ctx, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		key := &PaperKey{}
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

func (s *PostgresStore) Delete(ctx context.Context, ID string) error {
	_, err := s.stmtDelete.ExecContext(ctx, ID)
	return err
}

func (s *PostgresStore) DeleteByOwner(ctx context.Context, owner string) error {
	_, err := s.stmtDeleteByOwner.ExecContext(ctx, owner)
	return err
}
