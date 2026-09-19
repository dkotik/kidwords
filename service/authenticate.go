package service

import (
	"context"
	"errors"
	"fmt"
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

func (s *Service) AuthenticateUserWithNameAndPassword(
	ctx context.Context,
	name string,
	password string,
) (err error) {
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	userRepository, ok := rp.(secret.UserRepository)
	if !ok {
		return AuthenticationError{
			Cause: errors.ErrUnsupported,
		}
	}
	user, err := userRepository.RetrieveUserByName(ctx, name)
	if err != nil {
		return AuthenticationError{
			Cause: fmt.Errorf("unable to retrieve user: %w", err),
		}
	}
	return s.AuthenticateUserWithPassword(ctx, userRepository, user, password)
}

func (s *Service) AuthenticateUserWithEmailAndPassword(
	ctx context.Context,
	email string,
	password string,
) (err error) {
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	userRepository, ok := rp.(secret.UserRepository)
	if !ok {
		return AuthenticationError{
			Cause: errors.ErrUnsupported,
		}
	}
	user, err := userRepository.RetrieveUserByEmailAddress(ctx, email)
	if err != nil {
		return AuthenticationError{
			Cause: fmt.Errorf("unable to retrieve user: %w", err),
		}
	}
	return s.AuthenticateUserWithPassword(ctx, userRepository, user, password)
}

func (s *Service) AuthenticateUserWithPassword(
	ctx context.Context,
	rp secret.UserRepository,
	u User,
	password string,
) (err error) {
	hash, err := secret.ParseArgonHash(u.GetPasswordHash())
	if err != nil {
		return AuthenticationError{
			UserID: u.GetID(),
			Cause:  err,
		}
	}
	ok, err := hash.Match([]byte(password))
	if err != nil {
		return AuthenticationError{
			UserID: u.GetID(),
			Cause:  err,
		}
	}
	if !ok {
		return AuthenticationError{
			UserID: u.GetID(),
			Cause:  errors.New("password does not match"),
		}
	}
	if err = rp.MarkUserAsActive(ctx, u.GetID()); err != nil {
		return AuthenticationError{
			UserID: u.GetID(),
			Cause:  err,
		}
	}
	return nil
}

func (s *Service) AuthenticateUserWithNameAndPaperKey(
	ctx context.Context,
	name string,
	key string,
) (err error) {
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	userRepository, ok := rp.(secret.UserRepository)
	if !ok {
		return AuthenticationError{
			Cause: errors.ErrUnsupported,
		}
	}
	user, err := userRepository.RetrieveUserByName(ctx, name)
	if err != nil {
		return AuthenticationError{
			Cause: fmt.Errorf("unable to retrieve user: %w", err),
		}
	}
	return s.AuthenticateUserWithPaperKey(ctx, rp, user.GetID(), key)
}

func (s *Service) AuthenticateUserWithEmailAndPaperKey(
	ctx context.Context,
	email string,
	key string,
) (err error) {
	rp, tx, err := s.repository.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Close(&err)
	userRepository, ok := rp.(secret.UserRepository)
	if !ok {
		return AuthenticationError{
			Cause: errors.ErrUnsupported,
		}
	}
	user, err := userRepository.RetrieveUserByEmailAddress(ctx, email)
	if err != nil {
		return AuthenticationError{
			Cause: fmt.Errorf("unable to retrieve user: %w", err),
		}
	}
	return s.AuthenticateUserWithPaperKey(ctx, rp, user.GetID(), key)
}

func (s *Service) AuthenticateUserWithPaperKey(
	ctx context.Context,
	rp secret.Repository,
	userID string,
	key string,
) (err error) {
	shards, _, ok := s.decoder.Decode(key)
	if !ok {
		return AuthenticationError{
			UserID: userID,
			Cause:  fmt.Errorf("invalid paper key"),
		}
	}

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
