package mock

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"uuid"

	"github.com/dkotik/kidwords/internal"
	"github.com/dkotik/kidwords/service/secret"
)

var uuidSpace = uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")

type mock struct {
	mu          *sync.Mutex
	userSecrets map[string][]secret.Secret
}

func New() secret.Repository {
	return &mock{
		mu:          &sync.Mutex{},
		userSecrets: make(map[string][]secret.Secret, 1),
	}
}

func (m *mock) Create(_ context.Context, s secret.Secret) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s.ID = internal.NewDeterministicUUID(fmt.Sprintf("%d", len(m.userSecrets)))
	// s.ID = s.Name + s.Type
	// s.Name = s.ID
	s.Fingerprint = []byte(s.ID)
	// if s.ID == "" {
	// s.ID = uuid.New().String()
	// }
	userSecrets, _ := m.userSecrets[s.UserID]
	m.userSecrets[s.UserID] = append(userSecrets, s)
	return s.ID, nil
}

func (m *mock) Retrieve(_ context.Context, ID string) (secret.Secret, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, userSecrets := range m.userSecrets {
		for _, secret := range userSecrets {
			if secret.ID == ID {
				return secret, nil
			}
		}
	}
	return secret.Secret{}, secret.ErrNotFound
}

func (m *mock) Update(_ context.Context, s secret.Secret) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, userSecrets := range m.userSecrets {
		for i, secret := range userSecrets {
			if secret.ID == s.ID {
				userSecrets[i] = s
				return nil
			}
		}
	}
	return secret.ErrNotFound
}

func (m *mock) Delete(_ context.Context, ID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, userSecrets := range m.userSecrets {
		for j, secret := range userSecrets {
			if secret.ID == ID {
				// fmt.Println(secret.ID, ID)
				// panic("found")
				m.userSecrets[i] = slices.Delete(userSecrets, j, j+1)
				return nil
			}
		}
	}
	return nil // secret.ErrNotFound
}

func (m *mock) List(_ context.Context, userID string) ([]secret.Secret, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, userSecrets := range m.userSecrets {
		if id == userID {
			return userSecrets, nil
		}
	}
	return nil, nil
}
