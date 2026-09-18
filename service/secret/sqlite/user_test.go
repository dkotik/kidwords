package sqlite

import (
	"testing"
	"time"

	"github.com/dkotik/kidwords/service/secret"
	"zombiezen.com/go/sqlite"
)

func TestUserRepositoryInterfaceContract(t *testing.T) {
	conn, err := sqlite.OpenConn(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	rp1, err := New(conn)
	if err != nil {
		t.Fatal(err)
	}

	rp := rp1.(*sqRepository)

	u := &User{
		// ID:           "111",
		Name:         "testUser",
		PasswordHash: "????",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ActiveAt:     time.Now(),
	}
	ctx := t.Context()
	if err = rp.CreateUser(ctx, u); err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := rp.DeleteUser(ctx, u.ID)
		if err != nil {
			t.Fatal(err)
		}
	}()
	t.Log("user ID:", u.ID)
	existing, err := rp.RetrieveUser(ctx, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if existing == nil {
		t.Fatal("nil user")
	}
	if existing.ID != u.ID {
		t.Log("name:", existing.Name)
		t.Fatalf("expected user ID %q, got %q", u.ID, existing.ID)
	}

	secret.NewUserRepositoryTest(rp, u)(t)
}
