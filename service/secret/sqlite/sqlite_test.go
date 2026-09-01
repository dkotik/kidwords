package sqlite

import (
	"testing"

	"github.com/dkotik/kidwords/service/secret"
	"zombiezen.com/go/sqlite"
)

func TestRepositoryInterfaceContract(t *testing.T) {
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

	rp, err := New(conn)
	if err != nil {
		t.Fatal(err)
	}
	secret.NewSecretsRepositoryTest(rp)(t)
}
