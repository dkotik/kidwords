package service

import (
	"context"
	"testing"
	"time"
)

func NewSecretsRepositoryTest(r SecretRepository) func(*testing.T) {
	return func(t *testing.T) {
		userID := "testUser"
		secret, err := NewArgonSecret("test", "Test.", "password")
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = r.CreateSecret(ctx, userID, secret); err != nil {
			t.Fatal(err)
		}
		if err = r.UpdateSecret(ctx, userID, secret.ID, "test1", "Test1."); err != nil {
			t.Fatal(err)
		}

		all, err := r.ListSecrets(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
		if len(all) != 1 {
			t.Fatal("unexpected number of secrets:", len(all))
		}
		if all[0].Name != "test1" {
			t.Fatal("name does not match")
		}
		if all[0].Description != "Test1." {
			t.Fatal("description does not match")
		}
		if err = r.DeleteSecret(ctx, userID, secret.ID); err != nil {
			t.Fatal(err)
		}
		all, err = r.ListSecrets(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
		if len(all) != 0 {
			t.Fatal("failed to delete record")
		}
	}
}
