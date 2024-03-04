package htservice

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	// HashRecordLengthMax = 512
	HashArgonKeyLength  = 128
	HashArgonTimeCost   = 1
	HashArgonMemoryCost = 64 * 1024 // recommended by x/crypto/argon2
	HashArgonThreads    = 4
)

type ArgonSecret struct {
	ID              string
	Type            string // argon2id
	Name            string
	Description     string
	Secret          []byte
	Salt            []byte
	TimeCost        uint32
	MemoryCost      uint32
	ParallelStreams uint8
	Version         uint8 // 19
	Created         time.Time
}

func (a *ArgonSecret) Label() ArgonSecretLabel {
	return ArgonSecretLabel{
		ID:          a.ID,
		Type:        a.Type,
		Name:        a.Name,
		Description: a.Description,
		Version:     a.Version,
		Created:     a.Created,
	}
}

type ArgonSecretLabel struct {
	ID          string
	Type        string
	Name        string
	Description string
	Version     uint8
	Created     time.Time
}

func NewArgonSecret(name, description, secret string) (*ArgonSecret, error) {
	salt := make([]byte, HashArgonKeyLength)
	n, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	if n != HashArgonKeyLength {
		return nil, fmt.Errorf("not enough secure random bytes: %d vs %d", n, HashArgonKeyLength)
	}
	return &ArgonSecret{
		ID:          uuid.New().String(),
		Type:        "argon2id",
		Name:        name,
		Description: description,
		Secret: argon2.IDKey(
			[]byte(secret),
			salt,
			HashArgonTimeCost,
			HashArgonMemoryCost,
			HashArgonThreads,
			HashArgonKeyLength),
		Salt:            salt,
		TimeCost:        HashArgonTimeCost,
		MemoryCost:      HashArgonMemoryCost,
		ParallelStreams: HashArgonThreads,
		Version:         19,
		Created:         time.Now().UTC(),
	}, nil
}

func (a *ArgonSecret) Match(secret string) bool {
	return bytes.Equal(argon2.IDKey(
		[]byte(secret),
		a.Salt,
		a.TimeCost,
		a.MemoryCost,
		a.ParallelStreams,
		HashArgonKeyLength),
		a.Secret,
	)
}

func (a *ArgonSecret) String() string {
	return fmt.Sprintf(`$%s$v=%d$m=%d,t=%d,p=%d$%s$%s`,
		a.Type,
		a.Version,
		a.MemoryCost,
		a.TimeCost,
		a.ParallelStreams,
		base64.RawStdEncoding.EncodeToString(a.Salt),
		base64.RawStdEncoding.EncodeToString(a.Secret),
	)
}

type SecretRepository interface {
	CreateSecret(context.Context, string, *ArgonSecret) error
	// RetrieveSecret(context.Context, string, string) (*ArgonSecret, error)
	UpdateSecret(ctx context.Context, userID, id, name, description string) error
	DeleteSecret(context.Context, string, string) error
	ListSecrets(context.Context, string) ([]ArgonSecret, error)
}

type keyValueSecretRepository struct {
	kv KeyValueRepository
	// limit int
}

func NewKeyValueSecretRepository(kv KeyValueRepository) SecretRepository {
	return &keyValueSecretRepository{kv}
}

func (r *keyValueSecretRepository) decode(b []byte) (secrets []ArgonSecret, err error) {
	if len(b) < 1 {
		return nil, nil
	}
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(&secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}

func (r *keyValueSecretRepository) encode(secrets []ArgonSecret) ([]byte, error) {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(&secrets); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (r *keyValueSecretRepository) CreateSecret(
	ctx context.Context, userID string, secret *ArgonSecret,
) error {
	if secret == nil {
		return errors.New("cannot store a <nil> secret")
	}
	secrets, err := r.ListSecrets(ctx, userID)
	if err != nil {
		return err
	}
	if len(secrets) == 0 {
		b, err := r.encode([]ArgonSecret{*secret})
		if err != nil {
			return err
		}
		return r.kv.Set(ctx, []byte(userID), b)
	}

	return r.kv.Update(ctx, []byte(userID), func(b []byte) ([]byte, error) {
		secrets, err := r.decode(b)
		if err != nil {
			return nil, err
		}
		return r.encode(append(secrets, *secret))
	})
}

func (r *keyValueSecretRepository) UpdateSecret(
	ctx context.Context, userID, id, name, description string,
) error {
	return r.kv.Update(ctx, []byte(userID), func(b []byte) ([]byte, error) {
		secrets, err := r.decode(b)
		if err != nil {
			return nil, err
		}

		for i, secret := range secrets {
			if secret.ID == id {
				secret.Name = name
				secret.Description = description
				secrets[i] = secret
				return r.encode(secrets)
			}
		}

		return nil, fmt.Errorf("secret %q does not exist", id)
	})
}

func (r *keyValueSecretRepository) DeleteSecret(
	ctx context.Context, userID, id string,
) error {
	return r.kv.Update(ctx, []byte(userID), func(b []byte) ([]byte, error) {
		secrets, err := r.decode(b)
		if err != nil {
			return nil, err
		}

		for i, secret := range secrets {
			if secret.ID == id {
				secrets = append(secrets[:i], secrets[i+1:]...)
				return r.encode(secrets)
			}
		}

		return nil, fmt.Errorf("secret %q does not exist", id)
	})
}

func (r *keyValueSecretRepository) ListSecrets(
	ctx context.Context, userID string,
) ([]ArgonSecret, error) {
	b, err := r.kv.Get(ctx, []byte(userID))
	if err != nil {
		return nil, err
	}
	return r.decode(b)
}
