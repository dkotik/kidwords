package secret

import (
	"log/slog"
	"time"
)

type Secret struct {
	ID             string
	UserID         string
	Name           string
	Type           string
	Hash           string
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
