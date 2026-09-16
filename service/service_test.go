package service

import (
	"context"
	"testing"

	"github.com/dkotik/kidwords/service/secret"
	"github.com/dkotik/kidwords/service/secret/mock"
)

type mockUser struct{}

const mockUserID = "mockUserID"

func (u mockUser) GetID() string {
	return mockUserID
}

func (u mockUser) GetName() string {
	return "mockUserName"
}

type mockAuthenticator struct{}

func (m mockAuthenticator) Authenticate(context.Context) (User, error) {
	return mockUser{}, nil
}

func newTestService(t testing.TB, withOptions ...Option) (*Service, secret.Repository) {
	repository := mock.New()
	service, err := New(
		mockAuthenticator{},
		repository,
	)
	if err != nil {
		t.Fatal(err)
	}
	return service, repository
}
