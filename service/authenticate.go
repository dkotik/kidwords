package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/dkotik/kidwords"
	"github.com/dkotik/kidwords/service/secret"
)

type AuthenticationError struct {
	UserID string
	Cause  error
}

func (err AuthenticationError) Error() string {
	return "credentials did not match any known user"
}

func (err AuthenticationError) Unwrap() error {
	return err.Cause
}

func (err AuthenticationError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Bool("authenticated", false),
		slog.String("user_id", err.UserID),
		slog.Any("error", err.Cause),
	)
}

func (s *Service) Authenticate(
	ctx context.Context,
	rp secret.Repository,
	userID string,
	key string,
) (err error) {
	shards, err := s.decoder.Decode(key)
	if err != nil {
		return err
	}

	// rp, tx, err := s.repository.BeginTransaction(ctx)
	// if err != nil {
	// 	return fmt.Errorf("unable to begin transaction: %w", err)
	// }
	// defer tx.Close(&err)
	keys, err := rp.List(ctx, userID)
	if err != nil {
		return err
	}
	matched, ok := matchKey(keys, shards)
	if !ok {
		return AuthenticationError{
			UserID: userID,
			Cause:  err,
		}
	}
	hash, err := secret.ParseArgonHash(matched.SaltedHash)
	if err != nil {
		return AuthenticationError{
			UserID: userID,
			Cause:  err,
		}
	}
	secretBytes, err := kidwords.Combine(shards)
	if err != nil {
		return AuthenticationError{
			UserID: userID,
			Cause:  err,
		}
	}
	ok, err = hash.Match(secretBytes)
	if err != nil {
		return AuthenticationError{
			UserID: userID,
			Cause:  err,
		}
	}
	if !ok {
		return AuthenticationError{
			UserID: userID,
			Cause:  errors.New("recovered secret does not match"),
		}
	}
	matched.LastAcceptedAt = time.Now()
	if err = s.repository.Update(ctx, matched); err != nil {
		return AuthenticationError{
			UserID: userID,
			Cause:  err,
		}
	}
	return nil
}

func matchKey(keys []secret.Secret, shards []kidwords.Shard) (secret.Secret, bool) {
	dataLength := 0
	index := 0

nextKey:
	for _, key := range keys {
		for _, shard := range shards {
			dataLength = len(shard.Data)
			index = int(shard.Index) - 1
			// fmt.Printf("shard: %d, index: %d vs %d\n", shard.Index, key.Fingerprint[index], shard.Data[dataLength-1])
			if dataLength == 0 ||
				index >= len(key.Fingerprint) ||
				key.Fingerprint[index] != shard.Data[dataLength-1] {
				continue nextKey
			}
		}
		return key, true
	}
	return secret.Secret{}, false
}
