package test

import (
	"context"
	"testing"
	"time"

	"github.com/dkotik/kidwords/store"
	"github.com/google/uuid"
)

func RunAllOperations(s store.Store) func(*testing.T) {
	return func(t *testing.T) {
		hash, err := store.NewArgonHash([]byte("secret key"))
		if err != nil {
			t.Fatal("failed to construct the hash", hash)
		}

		paperKey := &store.PaperKey{
			ID:         uuid.New().String(),
			Owner:      uuid.New().String(),
			Name:       "Test Key",
			SaltedHash: hash.String(),
			Created:    time.Now().UTC(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err = s.Create(ctx, paperKey); err != nil {
			t.Fatal("failed to create paper key:", err)
		}

		keys, err := s.RetrieveAll(ctx, paperKey.Owner)
		if err != nil {
			t.Fatal("cannot retrieve keys:", err)
		}
		if keys[0].ID != paperKey.ID {
			t.Fatal("mismatch ID:", keys[0].ID, paperKey.ID)
		}
		if keys[0].Owner != paperKey.Owner {
			t.Fatal("mismatch owner:", keys[0].Owner, paperKey.Owner)
		}
		if keys[0].Name != paperKey.Name {
			t.Fatal("mismatch name:", keys[0].Name, paperKey.Name)
		}
		if keys[0].SaltedHash != paperKey.SaltedHash {
			t.Fatal("mismatch salted hash:", keys[0].SaltedHash, paperKey.SaltedHash)
		}
		if keys[0].Created.Unix() != paperKey.Created.Unix() {
			t.Fatal("mismatch creation date:", keys[0].Created, paperKey.Created)
		}

		if err = s.Delete(ctx, paperKey.ID); err != nil {
			t.Fatal("failed to create paper key:", err)
		}
	}
}
