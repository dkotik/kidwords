package secret

import (
	"context"
	"testing"
	"time"

	"github.com/dkotik/kidwords/service"
)

func NewSecretsRepositoryTest(r Repository) func(*testing.T) {
	return func(t *testing.T) {
		userID := "testUser"
		secret := service.Secret{
			Name: "test",
		}
		var err error
		// secret, err := NewArgonSecret("test", "Test.", "password")
		// if err != nil {
		// 	t.Fatal(err)
		// }
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = r.Create(ctx, secret); err != nil {
			t.Fatal(err)
		}
		if err = r.Update(ctx, secret); err != nil {
			t.Fatal(err)
		}

		all, err := r.List(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
		if len(all) != 1 {
			t.Fatal("unexpected number of secrets:", len(all))
		}
		if all[0].Name != "test1" {
			t.Fatal("name does not match")
		}
		if err = r.Delete(ctx, userID); err != nil {
			t.Fatal(err)
		}
		all, err = r.List(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
		if len(all) != 0 {
			t.Fatal("failed to delete record")
		}
	}
}
