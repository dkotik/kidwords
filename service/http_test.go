package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/kidwords/service/secret/mock"
)

type pageTester func(*http.Request) ([]byte, int, error)

type mockUser struct{}

func (u mockUser) GetID() string {
	return "mockUserID"
}

func (u mockUser) GetName() string {
	return "mockUserName"
}

type mockAuthenticator struct{}

func (m mockAuthenticator) Authenticate(context.Context) (User, error) {
	return mockUser{}, nil
}

func TestHandlers(t *testing.T) {
	prefix := "/"
	service, err := New(
		mockAuthenticator{},
		mock.New(),
	)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	if err = service.MountMux(
		mux,
		htadaptor.New(),
		prefix,
		nil,
	); err != nil {
		t.Fatal(err)
	}
	server := newMockServer(t, mux)
	if server == nil {
		t.Fatal("nil server")
	}

	t.Run("static assets", func(t *testing.T) {
		for _, path := range []string{
			"htmx.min.js",
			"bulma.min.css",
		} { // static assets should return 200 OK
			req, err := http.NewRequest("GET", prefix+path, nil)
			if err != nil {
				t.Fatal(err)
			}
			data, sc, err := server(req)
			if err != nil {
				t.Fatal(err)
			}
			if sc != http.StatusOK {
				t.Fatalf("%s: expected status code %d, got %d", path, http.StatusOK, sc)
			}
			if len(data) == 0 {
				t.Fatalf("%s: expected non-empty body", path)
			}
		}
	})

}

func newMockServer(t testing.TB, h http.Handler) pageTester {
	server := httptest.NewTestServer(t, h)
	t.Cleanup(server.Close)
	client := server.Client()
	prefix, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	return func(r *http.Request) (body []byte, statusCode int, err error) {
		cp := prefix.Clone()
		cp.Path = r.URL.Path
		cp.RawPath = r.URL.RawPath
		r.URL = cp
		// r.URL.Path = prefix.JoinPath(r.URL.Path).String()
		// r.URL.RawPath = prefix.JoinPath(r.URL.RawPath).String()
		resp, err := client.Do(r)
		if err != nil {
			return nil, 0, err
		}
		defer func() {
			err = errors.Join(err, resp.Body.Close())
		}()
		statusCode = resp.StatusCode
		body, err = io.ReadAll(resp.Body)
		return body, statusCode, err
	}
}
