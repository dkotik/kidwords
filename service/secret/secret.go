package secret

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const (
	TypePaperKey = "kidwordsPaperKey"
)

var (
	ErrDuplicateSecretName = errors.New("duplicate secret name")
	ErrNotFound            = errors.New("secret not found")
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
		return fmt.Errorf("ID mismatch: %q vs %q", s.ID, other.ID)
	}
	if s.UserID != other.UserID {
		return fmt.Errorf("user ID mismatch: %q vs %q", s.UserID, other.UserID)
	}
	if s.Name != other.Name {
		return fmt.Errorf("name is different: %q vs %q", s.Name, other.Name)
	}
	if s.Type != other.Type {
		return fmt.Errorf("type is different: %q vs %q", s.Type, other.Type)
	}
	if s.SaltedHash != other.SaltedHash {
		return fmt.Errorf("salted hash is different: %q vs %q", s.SaltedHash, other.SaltedHash)
	}
	if s.CreatedAt.Unix() != other.CreatedAt.Unix() {
		return fmt.Errorf("created at is different: %q vs %q", s.CreatedAt, other.CreatedAt)
	}
	if s.UpdatedAt.Unix() != other.UpdatedAt.Unix() {
		return fmt.Errorf("updated at is different: %q vs %q", s.UpdatedAt, other.UpdatedAt)
	}
	if s.LastAcceptedAt.Unix() != other.LastAcceptedAt.Unix() {
		return fmt.Errorf("last accepted at is different: %q vs %q", s.LastAcceptedAt, other.LastAcceptedAt)
	}
	return nil
}
