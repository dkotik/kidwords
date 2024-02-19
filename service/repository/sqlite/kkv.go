package sqlite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zombiezen.com/go/sqlite/sqlitex"
)

type KeyKeyValueRepository struct {
	dbpool    *sqlitex.Pool
	sqlGet    string
	sqlSet    string
	sqlUpdate string
	sqlDelete string
}

func NewKeyKeyValueRepository(withOptions ...Option) (_ *KeyKeyValueRepository, err error) {
	o := &options{}
	for _, option := range append(
		withOptions, WithDefaultTableName(), WithDefaultConnection(),
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("unable to create key-key-value Sqlite3 repository: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	dbpool, err := sqlitex.NewPool(o.Connection, sqlitex.PoolOptions{
		PoolSize: o.ConnectionCount,
		// PrepareConn: sqlitex.ConnPrepareFunc(func(conn *sqlite.Conn) error {
		// 	return nil
		// }),
	})
	if err != nil {
		return nil, err
	}

	conn := dbpool.Get(ctx)
	if conn == nil {
		return nil, errors.New("no available connection")
	}
	defer dbpool.Put(conn)

	err = sqlitex.ExecuteTransient(conn, fmt.Sprintf(`
	  CREATE TABLE IF NOT EXISTS %s (
	    first BLOB,
	    second BLOB,
	    value BLOB,
      PRIMARY KEY (first, second)
	  );`, o.Table), nil)
	if err != nil {
		return nil, err
	}

	return &KeyKeyValueRepository{
		dbpool:    dbpool,
		sqlGet:    fmt.Sprintf("SELECT value FROM %s WHERE first=$first AND second=$second;", o.Table),
		sqlSet:    fmt.Sprintf("INSERT INTO %s (first, second, value) VALUES($first, $second, $value);", o.Table),
		sqlUpdate: fmt.Sprintf("UPDATE %s SET value=$value WHERE first=$first AND second=$second;", o.Table),
		sqlDelete: fmt.Sprintf("DELETE FROM %s WHERE first=$first AND second=$second;", o.Table),
	}, nil
}

func (kkv *KeyKeyValueRepository) Get(ctx context.Context, first, second []byte) (result []byte, err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return nil, errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlGet)
	stmt.BindBytes(1, first)
	stmt.BindBytes(2, second)

	hasRow, err := stmt.Step()
	if err != nil {
		return nil, err
	}
	if hasRow {
		if l := stmt.GetLen("value"); l > 0 {
			result = make([]byte, l)
			stmt.GetBytes("value", result)
		}
	}
	if err = stmt.Reset(); err != nil {
		return nil, err
	}
	return result, nil
}

func (kkv *KeyKeyValueRepository) Set(ctx context.Context, first, second, value []byte) (err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlSet)
	stmt.BindBytes(1, first)
	stmt.BindBytes(2, second)
	stmt.BindBytes(3, value)

	_, err = stmt.Step()
	if err != nil {
		return err
	}
	return stmt.Reset()
}

func (kkv *KeyKeyValueRepository) Update(ctx context.Context, first, second []byte, update func([]byte) ([]byte, error)) (err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	txClose, err := sqlitex.ExclusiveTransaction(conn)
	if err != nil {
		return err
	}
	defer txClose(&err)

	stmt := conn.Prep(kkv.sqlGet)
	stmt.BindBytes(1, first)
	stmt.BindBytes(2, second)
	hasRows, err := stmt.Step()
	if err != nil {
		return err
	}
	if !hasRows {
		return fmt.Errorf("value %s/%s does not exist", string(first), string(second))
	}

	var value []byte
	if l := stmt.GetLen("value"); l > 0 {
		value = make([]byte, l)
		stmt.GetBytes("value", value)
	}
	if err = stmt.Reset(); err != nil {
		return err
	}

	value, err = update(value)
	if err != nil {
		return err
	}

	stmt = conn.Prep(kkv.sqlUpdate)
	stmt.BindBytes(1, value) // mind the order!
	stmt.BindBytes(2, first)
	stmt.BindBytes(3, second)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return stmt.Reset()
}

func (kkv *KeyKeyValueRepository) Delete(ctx context.Context, first, second []byte) (err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlDelete)
	stmt.SetBytes("$first", first)
	stmt.SetBytes("$second", second)

	_, err = stmt.Step()
	return err
}
