package secret

import (
	"context"
	"testing"
	"time"
)

func NewSecretsRepositoryTest(r Repository) func(*testing.T) {
	return func(t *testing.T) {
		now := time.Now()
		secret := Secret{
			Type:      TypePaperKey,
			Name:      "test",
			CreatedAt: now,
			UpdatedAt: now,
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
		if s2.ID == "" {
			t.Fatal("ID is empty")
		}
		if err = secret.IsEqual(s2); err != nil {
			t.Fatal(err)
		}
		// if err = r.Update(ctx, secret); err != nil {
		// 	t.Fatal(err)
		// }

		// all, err := r.List(ctx, userID)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// if len(all) != 1 {
		// 	t.Fatal("unexpected number of secrets:", len(all))
		// }
		// if all[0].Name != "test1" {
		// 	t.Fatal("name does not match")
		// }
		// if err = r.Delete(ctx, userID); err != nil {
		// 	t.Fatal(err)
		// }
		// all, err = r.List(ctx, userID)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// if len(all) != 0 {
		// 	t.Fatal("failed to delete record")
		// }
	}
}
