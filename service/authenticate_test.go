package service

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestAuthentication(t *testing.T) {
	service, repository := newTestService(t)
	ctx := t.Context()
	secret, err := service.CreateKey(ctx, "mockKey")
	if err != nil {
		t.Fatal(err)
	}
	user, err := service.authenticator.Authenticate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	for _, row := range secret.Secret {
		for _, cell := range row {
			_, _ = fmt.Fprintf(b, "%d", cell.Index)
			_ = b.WriteByte(' ')
			for _, word := range cell.Words {
				_, _ = b.WriteString(word)
				_ = b.WriteByte(' ')
			}
		}
	}

	err = service.Authenticate(
		ctx,
		repository,
		user.GetID(),
		b.String(),
	)
	if err != nil {
		asAuth, ok := errors.AsType[AuthenticationError](err)
		if ok {
			err = asAuth.Cause
		}
		t.Fatal(err)
	}
}
