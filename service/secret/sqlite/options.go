package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dkotik/kidwords/service/secret"
)

const (
	DefaultUsersTableName   = `users`
	DefaultSecretsTableName = `kidwords_secrets`
)

type options struct {
	UsersTableName   string
	SecretsTableName string
	SecretType       string
	UserForeignKey   *ForeignKey
}

func (o *options) getInstallScript() string {
	b := &strings.Builder{}
	if o.UserForeignKey == nil {
		_, _ = b.WriteString(fmt.Sprintf(`
		    CREATE TABLE IF NOT EXISTS %s (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					password_hash TEXT NOT NULL,
					email TEXT NOT NULL UNIQUE,
					email_verified INTEGER DEFAULT 0
						CHECK (email_verified IN (0, 1)),
					created_at TEXT,
					updated_at TEXT,
					active_at TEXT
		    );
		  `,
			escapeIdentifier(o.UsersTableName),
		))
		o.UserForeignKey = &ForeignKey{
			TableName: o.UsersTableName,
			FieldName: "id",
		}
	}
	_, _ = b.WriteString(fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
      %s,
      FOREIGN KEY (user_id) REFERENCES %s(%s),
      UNIQUE (user_id, name),
      UNIQUE (user_id, fingerprint)
    );`,
		escapeIdentifier(o.SecretsTableName),
		SecretsTableFields,
		escapeIdentifier(o.UserForeignKey.TableName),
		escapeIdentifier(o.UserForeignKey.FieldName),
	))
	return b.String()
}

type ForeignKey struct {
	TableName string
	FieldName string
}

type Option func(*options) error

func WithUsersTableName(name string) Option {
	return func(o *options) error {
		if name == "" {
			return errors.New("cannot use an empty table name")
		}
		if o.UsersTableName != "" {
			return errors.New("users table name is already set")
		}
		if o.UserForeignKey != nil {
			return errors.New("choose either users table name or users foreign key; cannot set both")
		}
		o.UsersTableName = name
		return nil
	}
}

func WithSecretsTableName(name string) Option {
	return func(o *options) error {
		if name == "" {
			return errors.New("cannot use an empty table name")
		}
		if o.SecretsTableName != "" {
			return errors.New("secrets table name is already set")
		}
		o.SecretsTableName = name
		return nil
	}
}

func WithSecretType(secretType string) Option {
	return func(o *options) error {
		if secretType == "" {
			return errors.New("cannot use an empty secret type")
		}
		if o.SecretType != "" {
			return errors.New("secret type is already set")
		}
		o.SecretType = secretType
		return nil
	}
}

func WithUserForeignKey(tableName, fieldName string) Option {
	return func(o *options) error {
		tableName = strings.TrimSpace(tableName)
		if tableName == "" {
			return errors.New("table name is required")
		}
		fieldName = strings.TrimSpace(fieldName)
		if fieldName == "" {
			return errors.New("field name is required")
		}
		if o.UsersTableName != "" {
			return errors.New("choose either users table name or users foreign key; cannot set both")
		}
		if o.UserForeignKey != nil {
			return errors.New("foreign key is already set")
		}
		o.UserForeignKey = &ForeignKey{
			TableName: tableName,
			FieldName: fieldName,
		}
		return nil
	}
}

func withDefaultTableNames() Option {
	return func(o *options) (err error) {
		if o.UsersTableName == "" {
			err = WithUsersTableName(DefaultUsersTableName)(o)
			if err != nil {
				return err
			}
		}

		if o.SecretsTableName == "" {
			err = WithSecretsTableName(DefaultSecretsTableName)(o)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func withDefaultSecretType() Option {
	return func(o *options) error {
		if o.SecretType != "" {
			return nil
		}
		return WithSecretType(secret.TypePaperKey)(o)
	}
}
