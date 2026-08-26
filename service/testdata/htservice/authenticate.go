package htservice

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

type Authenticator interface {
	UserID(context.Context) (string, error)
	SessionID(context.Context) (string, error)
	IsAdministrator(context.Context) bool
}

type AuthenticationAttempt struct {
	ID        string
	SecretID  string
	SessionID string
	Error     string
	Address   string
	Time      time.Time
}

type AuthenticationAttemptRepository interface {
	CreateAuthenticationAttempt(context.Context, string, AuthenticationAttempt) error
	RecentAuthenticationAttempts(context.Context, string, time.Duration) ([]AuthenticationAttempt, error)
}

type keyValueAuthenticationAttemptRepository struct {
	kv KeyValueRepository
}

func (r *keyValueAuthenticationAttemptRepository) CreateAuthenticationAttempt(ctx context.Context, userID string, a AuthenticationAttempt) error {
	key := []byte(userID)
	b, err := r.kv.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		buf := &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode([]AuthenticationAttempt{a}); err != nil {
			return err
		}
		return r.kv.Set(ctx, key, buf.Bytes())
	}

	return r.kv.Update(ctx, key, func(b []byte) (_ []byte, err error) {
		var attempts []AuthenticationAttempt
		if err = json.NewDecoder(bytes.NewReader(b)).Decode(&attempts); err != nil {
			return nil, err
		}
		buf := &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode(append(attempts, a)); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	})
}

func (r *keyValueAuthenticationAttemptRepository) RecentAuthenticationAttempts(ctx context.Context, userID string, d time.Duration) (result []AuthenticationAttempt, err error) {
	b, err := r.kv.Get(ctx, []byte(userID))
	if err != nil {
		return nil, err
	}
	var attempts []AuthenticationAttempt
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(&attempts); err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-d)
	for _, attempt := range attempts {
		if attempt.Time.After(cutoff) {
			result = append(result, attempt)
		}
	}
	return
}
