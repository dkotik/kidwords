package sqlite

import "errors"

type options struct {
	Connection      string
	Table           string
	ConnectionCount int
}

type Option func(*options) error

func WithConnection(dns string) Option {
	if dns == ":memory:" {
		return WithMemoryConnection()
	}

	return func(o *options) error {
		if o.Connection != "" {
			return errors.New("connection DNS is already set")
		}
		if dns == "" {
			return errors.New("cannot use an empty connection DNS")
		}
		o.Connection = dns
		return nil
	}
}

func WithMemoryConnection() Option {
	return func(o *options) (err error) {
		o.ConnectionCount = 1
		return WithConnection("file:memory:?mode=memory")(o)
	}
}

func WithDefaultConnection() Option {
	return func(o *options) (err error) {
		if o.Connection != "" {
			return nil
		}
		return WithMemoryConnection()(o)
	}
}

func WithConnectionCount(limit int) Option {
	return func(o *options) error {
		if limit < 1 {
			return errors.New("invalid connection count")
		}
		if o.ConnectionCount != 0 {
			return errors.New("connection count is already set")
		}
		o.ConnectionCount = limit
		return nil
	}
}

func WithDefaultConnectionCountOf20() Option {
	return func(o *options) error {
		if o.ConnectionCount > 1 {
			return nil
		}
		return WithConnectionCount(10)(o)
	}
}

func WithTable(name string) Option {
	return func(o *options) error {
		if o.Table != "" {
			return errors.New("table name is already set")
		}
		if name == "" {
			return errors.New("cannot use empty table name")
		}
		o.Table = name
		return nil
	}
}

func WithDefaultTableName() Option {
	return func(o *options) error {
		if o.Table != "" {
			return nil
		}
		return WithTable("key_key_value_repository")(o)
	}
}
