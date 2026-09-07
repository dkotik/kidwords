package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dkotik/kidwords/service/secret"
)

const DefaultTableName = `kidwords_paper_keys`

type options struct {
	TableName      string
	SecretType     string
	UserForeignKey *ForeignKey
}

func (o *options) getInstallScript() string {
	if o.UserForeignKey == nil {
		return fmt.Sprintf(`
      CREATE TABLE IF NOT EXISTS %s (
        %s
      );`,
			escapeIdentifier(o.TableName),
			SecretsTableFields,
		)
	}
	return fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
      %s,
      FOREIGN KEY (user_id) REFERENCES %s(%s),
      UNIQUE (user_id, fingerprint)
    );`,
		escapeIdentifier(o.TableName),
		SecretsTableFields,
		escapeIdentifier(o.UserForeignKey.TableName),
		escapeIdentifier(o.UserForeignKey.FieldName),
	)
}

type ForeignKey struct {
	TableName string
	FieldName string
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

func withDefaultTableName() Option {
	return func(o *options) error {
		if o.TableName != "" {
			return nil
		}
		return WithTableName(DefaultTableName)(o)
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
