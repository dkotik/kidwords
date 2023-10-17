package store

import "errors"

const DefaultTableName = `kidwords_paper_keys`

type options struct {
	TableName   string
	CreateTable bool
}

type Option func(*options) error

func WithTableName(name string) Option {
	return func(o *options) error {
		if name == "" {
			return errors.New("cannot use an empty table name")
		}
		if o.TableName != "" {
			return errors.New("table name is already set")
		}
		o.TableName = name
		return nil
	}
}

func WithGuaranteedTable() Option {
	return func(o *options) error {
		if o.CreateTable {
			return errors.New("table creation is already guaranteed")
		}
		o.CreateTable = true
		return nil
	}
}

func withDefaultTableName() Option {
	return func(o *options) error {
		if o.TableName != "" {
			return nil
		}
		return WithTableName(DefaultTableName)(o)
	}
}
