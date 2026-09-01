package sqlite

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dkotik/kidwords/service/secret"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const SecretsTableFields = `
      id                  TEXT PRIMARY KEY,
      user_id             TEXT NOT NULL,
      name                TEXT NOT NULL,
      salted_hash         TEXT NOT NULL,
      created_at          TEXT NOT NULL,
      updated_at          TEXT NOT NULL,
      last_accepted_at    TEXT`

var _ secret.Repository = (*sqRepository)(nil) // interface satisfaction

type sqRepository struct {
	Conn              *sqlite.Conn
	stmtCreate        *sqlite.Stmt
	stmtRetrieveAll   *sqlite.Stmt
	stmtDelete        *sqlite.Stmt
	stmtDeleteByOwner *sqlite.Stmt
}

func New(conn *sqlite.Conn, withOptions ...Option) (s *sqRepository, err error) {
	if conn == nil {
		return nil, errors.New("cannot use a <nil> database connection")
	}
	o := &options{}
	for _, option := range append(withOptions, withDefaultTableName()) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize the store: %w", err)
		}
	}

	if err = sqlitex.ExecScript(conn, o.getInstallScript()); err != nil {
		return nil, fmt.Errorf("cannot create the tables: %w", err)
	}

	s = &sqRepository{
		Conn: conn,
	}

	if s.stmtCreate, err = conn.Prepare(
		fmt.Sprintf(`
      INSERT INTO %s(id, owner, name, saltedHash, created) VALUES($1, $2, $3, $4, $5);
    `, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtRetrieveAll, err = conn.Prepare(
		fmt.Sprintf(`
      SELECT
        id, owner, name, saltedHash, created
      FROM %s WHERE owner=$1 ORDER BY created DESC;`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtDelete, err = conn.Prepare(
		fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtDeleteByOwner, err = conn.Prepare(
		fmt.Sprintf(`DELETE FROM %s WHERE owner=$1;`, o.TableName),
	); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *sqRepository) BindContext(ctx context.Context) func() {
	old := s.Conn.SetInterrupt(ctx.Done())
	return func() {
		s.Conn.SetInterrupt(old)
	}
}

func (s *sqRepository) BeginTransaction(context.Context) (secret.Repository, secret.Transaction, error) {
	return s, transaction(sqlitex.Transaction(s.Conn)), nil
}

func (s *sqRepository) WithTransaction(
	ctx context.Context,
	tx secret.Transaction,
) (secret.Repository, error) {
	_, ok := tx.(transaction)
	if !ok {
		return nil, fmt.Errorf("imcompatible transaction type")
	}
	return s, nil
}

// escapeIdentifier safely quotes an SQLite table or column name.
func escapeIdentifier(name string) string {
	// Double quotes are escaped by doubling them in SQL identifiers
	escaped := strings.ReplaceAll(name, `"`, `""`)
	return fmt.Sprintf(`"%s"`, escaped)
}

func encodeTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func decodeTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

type transaction func(*error)

func (t transaction) Close(err *error) {
	t(err)
}
