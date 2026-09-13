package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/service/secret"
)

type Secret struct {
	*KeyView
	Secret kidwords.Table
	Quorum int
}

func (s *Service) CreateKey(ctx context.Context, name string) (v *Secret, err error) {
	user, err := s.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate user: %w", err)
	}

	secretBytes := make([]byte, s.keyLength)
	if _, err = rand.Read(secretBytes); err != nil {
		return nil, err
	}
	argonHash, err := secret.NewArgonHash(secretBytes)
	if err != nil {
		return nil, err
	}

	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	userID := user.GetID()
	keys, err := rp.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(keys) >= s.keyCountLimit {
		return nil, errors.New("too many paper keys, please delete one")
	}

	kidwordsSecret, err := kidwords.NewSecret(
		secretBytes,
		s.shardCount,
		s.quorumCount,
	)
	if err != nil {
		return nil, err
	}
	v = &Secret{
		Secret: s.encoder.MakeTable(kidwordsSecret),
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = newName(v.Secret)
	}
	ID, err := rp.Create(ctx, secret.Secret{
		UserID:      userID,
		Name:        name,
		Fingerprint: kidwordsSecret.GetFingerprint(),
		Type:        secret.TypePaperKey,
		SaltedHash:  argonHash.String(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}

	scrt, err := s.repository.Retrieve(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve secret: %w", err)
	}
	v.KeyView, err = NewKeyView(user, scrt)
	if err != nil {
		return nil, fmt.Errorf("unable to create key view: %w", err)
	}
	return v, nil
}

func newName(s kidwords.Table) string {
	b := &strings.Builder{}
	i := 0
	for _, row := range s {
		for _, cell := range row {
			for _, word := range cell.Words {
				if i > 0 {
					_ = b.WriteByte(' ')
				}
				_, _ = b.WriteString(word)
				i++
				if i == 4 {
					return b.String()
				}
			}
		}
	}
	return b.String()
}
