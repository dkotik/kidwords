package htservice

import (
	"testing"

	"github.com/dkotik/kidwords/htservice/repository/sqlite"
)

func TestSecretRepository(t *testing.T) {
	kkv, err := sqlite.NewKeyKeyValueRepository()
	if err != nil {
		t.Fatal(err)
	}
	r := NewKeyValueSecretRepository(
		NewKeyValueFromKeyKeyValueRepository("testSecrets", kkv))

	t.Run("sqlite implementation", NewSecretsRepositoryTest(r))
}
