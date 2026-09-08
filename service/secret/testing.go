package secret

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

func NewRepositoryTest(r Repository) func(*testing.T) {
	return func(t *testing.T) {
		if r == nil {
			t.Fatal("nil repository")
		}

		now := time.Now()
		secret := Secret{
			UserID:      "user1",
			Type:        TypePaperKey,
			Name:        "test",
			Fingerprint: []byte(`12345`),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		var err error

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if secret.ID, err = r.Create(ctx, secret); err != nil {
			t.Fatal(err)
		}
		s2, err := r.Retrieve(ctx, secret.ID)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err = r.Delete(ctx, secret.ID); err != nil {
				t.Fatal("unable to delete secret:", err)
			}
			s2, err = r.Retrieve(ctx, secret.ID)
			if err == nil {
				t.Fatal("expected ErrSecretNotFound, got nil")
			}
			if !errors.Is(err, ErrNotFound) {
				t.Fatal(err)
			}
		}()

		if s2.ID == "" {
			t.Fatal("ID is empty")
		}
		if s2.UserID == "" {
			t.Fatal("UserID is empty")
		}
		if err = secret.IsEqual(s2); err != nil {
			t.Fatal(err)
		}
		// secret.UserID = secret.UserID + ":updated"
		secret.Name = secret.Name + ":updated"
		secret.Type = secret.Type + ":updated"
		secret.SaltedHash = secret.SaltedHash + ":updated"
		now = now.Add(time.Minute * 12)
		secret.CreatedAt = now
		// secret.UpdatedAt = now
		secret.LastAcceptedAt = now

		if err = r.Update(ctx, secret); err != nil {
			t.Fatal(err)
		}
		s2, err = r.Retrieve(ctx, secret.ID)
		if err != nil {
			t.Fatal(err)
		}
		if s2.ID == "" {
			t.Fatal("ID is empty")
		}
		if s2.UserID == "" {
			t.Fatal("UserID is empty")
		}
		s2.UpdatedAt = secret.UpdatedAt
		if err = secret.IsEqual(s2); err != nil {
			t.Fatal(err)
		}
		// s2 = secret

		////////// LIST
		s3 := s2
		s3.ID = ""
		s3.Name = "third secret"
		if s3.ID, err = r.Create(ctx, s3); err != nil {
			t.Fatal(err)
		}
		s3, err = r.Retrieve(ctx, secret.ID)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err = r.Delete(ctx, s3.ID); err != nil {
				t.Fatal(err)
			}
			s3, err = r.Retrieve(ctx, s3.ID)
			if err == nil {
				t.Fatal("expected ErrSecretNotFound, got nil")
			}
			if !errors.Is(err, ErrNotFound) {
				t.Fatal(err)
			}
		}()

		all, err := r.List(ctx, s3.UserID)
		if err != nil {
			t.Fatal(err)
		}
		if len(all) != 2 {
			t.Fatal("unexpected number of secrets:", len(all))
		}
		slices.SortFunc(all, func(a, b Secret) int {
			return strings.Compare(a.Name, b.Name)
		})
		s2.ID = all[0].ID
		if err = s2.IsEqual(all[0]); err != nil {
			t.Fatal(err)
		}
		// s3.ID = all[1].ID
		// if err = s3.IsEqual(all[1]); err != nil {
		// 	t.Fatal(err)
		// }
	}
}
