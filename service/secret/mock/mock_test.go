package mock

import (
	"testing"

	"github.com/dkotik/kidwords/service/secret"
)

func TestRepositoryInterfaceContract(t *testing.T) {
	secret.NewSecretsRepositoryTest(New())(t)
}
