package sqlite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"zombiezen.com/go/sqlite/sqlitex"
)

type EphemeralKeyKeyValueRepository struct {
	dbpool    *sqlitex.Pool
	timeout   time.Duration
	sqlGet    string
	sqlSet    string
	sqlUpdate string
	sqlDelete string
	sqlClean  string
}

func NewEphemeral(timeout time.Duration, withOptions ...Option) (_ *EphemeralKeyKeyValueRepository, err error) {
	if timeout < time.Millisecond {
		return nil, errors.New("timeout cannot be lower than a millisecond")
	}
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
      expiry INTEGER,
      PRIMARY KEY (first, second)
	  );`, o.Table), nil)
	if err != nil {
		return nil, err
	}

	kkv := &EphemeralKeyKeyValueRepository{
		dbpool:    dbpool,
		sqlGet:    fmt.Sprintf("SELECT value FROM %s WHERE first=$first AND second=$second AND expiry>=$cutoff;", o.Table),
		sqlSet:    fmt.Sprintf("INSERT INTO %s (first, second, value, expiry) VALUES($first, $second, $value, $expiry);", o.Table),
		sqlUpdate: fmt.Sprintf("UPDATE %s SET value=$value, expiry=$expiry WHERE first=$first AND second=$second;", o.Table),
		sqlDelete: fmt.Sprintf("DELETE FROM %s WHERE first=$first AND second=$second;", o.Table),
		sqlClean:  fmt.Sprintf("DELETE FROM %s WHERE expiry<$cutoff;", o.Table),
	}

	ticker := time.NewTicker(timeout / 2)
	go func(ctx context.Context, kkv *EphemeralKeyKeyValueRepository, c <-chan time.Time) {
		var err error
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-c:
				if err = kkv.Clean(ctx, t); err != nil {
					slog.Default().Error(
						"failed key-key-value clean up procedure",
						slog.Any("error", err),
					)
				}
			}
		}
	}(context.Background(), kkv, ticker.C)

	return kkv, nil
}

func (kkv *EphemeralKeyKeyValueRepository) expiry() int64 {
	return time.Now().Add(kkv.timeout).Unix()
}

func (kkv *EphemeralKeyKeyValueRepository) Get(ctx context.Context, first, second []byte) (result []byte, err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return nil, errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlGet)
	stmt.BindBytes(1, first)
	stmt.BindBytes(2, second)
	stmt.BindInt64(3, time.Now().Unix())

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

func (kkv *EphemeralKeyKeyValueRepository) Set(ctx context.Context, first, second, value []byte) (err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlSet)
	stmt.BindBytes(1, first)
	stmt.BindBytes(2, second)
	stmt.BindBytes(3, value)
	stmt.BindInt64(4, kkv.expiry())
	_, err = stmt.Step()
	if err != nil {
		return err
	}
	return stmt.Reset()
}

func (kkv *EphemeralKeyKeyValueRepository) Update(ctx context.Context, first, second []byte, update func([]byte) ([]byte, error)) (err error) {
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
	stmt.BindInt64(3, time.Now().Unix())
	hasRows, err := stmt.Step()
	if err != nil {
		return err
	}
	var value []byte
	if !hasRows {
		if err = stmt.Reset(); err != nil {
			return err
		}
		value, err = update(nil)
		if err != nil {
			return err
		}
		stmt = conn.Prep(kkv.sqlSet)
		stmt.BindBytes(1, first)
		stmt.BindBytes(2, second)
		stmt.BindBytes(3, value)
		stmt.BindInt64(4, kkv.expiry())
		_, err = stmt.Step()
		if err != nil {
			return err
		}
		return errors.Join(err, stmt.Reset())
	}

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
	if len(value) == 0 {
		stmt = conn.Prep(kkv.sqlDelete)
		stmt.BindBytes(1, first)
		stmt.BindBytes(2, second)
		_, err = stmt.Step()
		return errors.Join(err, stmt.Reset())
	}

	stmt = conn.Prep(kkv.sqlUpdate)
	stmt.BindBytes(1, value) // mind the order!
	stmt.BindInt64(2, kkv.expiry())
	stmt.BindBytes(3, first)
	stmt.BindBytes(4, second)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return stmt.Reset()
}

func (kkv *EphemeralKeyKeyValueRepository) Delete(ctx context.Context, first, second []byte) (err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlDelete)
	stmt.BindBytes(1, first)
	stmt.BindBytes(2, second)

	_, err = stmt.Step()
	return err
}

func (kkv *EphemeralKeyKeyValueRepository) Clean(ctx context.Context, upto time.Time) (err error) {
	conn := kkv.dbpool.Get(ctx)
	if conn == nil {
		return errors.New("no available connection")
	}
	defer kkv.dbpool.Put(conn)

	stmt := conn.Prep(kkv.sqlClean)
	stmt.BindInt64(1, upto.Unix())

	_, err = stmt.Step()
	return err
}
