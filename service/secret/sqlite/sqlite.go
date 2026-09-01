package sqlite

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dkotik/kidwords/service/secret"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const SecretsTableFields = `
      id                  TEXT PRIMARY KEY,
      user_id             TEXT NOT NULL,
      name                TEXT NOT NULL UNIQUE,
      type                TEXT NOT NULL,
      salted_hash         TEXT NOT NULL,
      created_at          TEXT NOT NULL,
      updated_at          TEXT NOT NULL,
      last_accepted_at    TEXT`

var _ secret.Repository = (*sqRepository)(nil) // interface satisfaction

type sqRepository struct {
	mu           *sync.Mutex
	conn         *sqlite.Conn
	stmtCreate   *sqlite.Stmt
	stmtRetrieve *sqlite.Stmt
	stmtUpdate   *sqlite.Stmt
	stmtDelete   *sqlite.Stmt
	stmtList     *sqlite.Stmt
}

func New(conn *sqlite.Conn, withOptions ...Option) (s *sqRepository, err error) {
	if conn == nil {
		return nil, errors.New("cannot use a <nil> database connection")
	}
	o := &options{}
	for _, option := range append(
		withOptions,
		withDefaultTableName(),
		withDefaultSecretType(),
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize the store: %w", err)
		}
	}

	if err = sqlitex.ExecScript(conn, o.getInstallScript()); err != nil {
		return nil, fmt.Errorf("cannot create the tables: %w", err)
	}

	s = &sqRepository{
		mu:   &sync.Mutex{},
		conn: conn,
	}

	o.TableName = escapeIdentifier(o.TableName)
	if s.stmtCreate, err = conn.Prepare(
		fmt.Sprintf(`
      INSERT INTO %s(id, user_id, name, type, salted_hash, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)
    `, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtRetrieve, err = conn.Prepare(
		fmt.Sprintf(`
      SELECT user_id, name, type, salted_hash, created_at, updated_at, last_accepted_at
      FROM %s WHERE id=?`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtUpdate, err = conn.Prepare(
		fmt.Sprintf(`
      UPDATE %s SET name=?, type=?, salted_hash=?, created_at=?, updated_at=?, last_accepted_at=? WHERE id=?`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtList, err = conn.Prepare(
		fmt.Sprintf(`
      SELECT id, name, type, salted_hash, created_at, updated_at, last_accepted_at
      FROM %s WHERE user_id=? OR 1`, o.TableName),
	); err != nil {
		return nil, err
	}

	if s.stmtDelete, err = conn.Prepare(
		fmt.Sprintf(`DELETE FROM %s WHERE id=?`, o.TableName),
	); err != nil {
		return nil, err
	}

	return s, nil
}

func (r *sqRepository) BindContext(ctx context.Context) func() {
	r.mu.Lock()
	old := r.conn.SetInterrupt(ctx.Done())
	return func() {
		r.conn.SetInterrupt(old)
		r.mu.Unlock()
	}
}

func (r *sqRepository) BeginTransaction(context.Context) (secret.Repository, secret.Transaction, error) {
	return r, transaction(sqlitex.Transaction(r.conn)), nil
}

func (r *sqRepository) WithTransaction(
	ctx context.Context,
	tx secret.Transaction,
) (secret.Repository, error) {
	_, ok := tx.(transaction)
	if !ok {
		return nil, fmt.Errorf("imcompatible transaction type")
	}
	return r, nil
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
