package secret

import (
	"errors"
	"log/slog"
	"time"
)

type Secret struct {
	ID             string
	UserID         string
	Name           string
	Type           string
	SaltedHash     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastAcceptedAt time.Time
}

func (p *Secret) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", p.ID),
		slog.String("name", p.Name),
		slog.String("type", p.Type),
	)
}

func (s Secret) IsEqual(other Secret) error {
	if s.ID != other.ID {
		return errors.New("ID mismatch")
	}
	if s.UserID != other.UserID {
		return errors.New("user ID mismatch")
	}
	if s.Name != other.Name {
		return errors.New("name is different")
	}
	if s.Type != other.Type {
		return errors.New("type is different")
	}
	if s.SaltedHash != other.SaltedHash {
		return errors.New("salted hash is different")
	}
	if s.CreatedAt.Unix() != other.CreatedAt.Unix() {
		return errors.New("created at is different")
	}
	if s.UpdatedAt.Unix() != other.UpdatedAt.Unix() {
		return errors.New("updated at is different")
	}
	if s.LastAcceptedAt.Unix() != other.LastAcceptedAt.Unix() {
		return errors.New("last accepted at is different")
	}
	return nil
}
