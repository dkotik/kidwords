package service

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

type pageTester func(*http.Request) ([]byte, int, error)

func TestHandlers(t *testing.T) {
	// service := New(
	// 	authenticator Authenticator,
	// 	repository secret.Repository,
	// )
	// server := newMockServer(t, h http.Handler)
}

func newMockServer(t testing.TB, h http.Handler) pageTester {
	server := httptest.NewTestServer(t, h)
	t.Cleanup(server.Close)
	client := server.Client()
	prefix := server.URL
	return func(r *http.Request) (body []byte, statusCode int, err error) {
		r.URL.Path = path.Join(prefix, r.URL.Path)
		r.URL.RawPath = path.Join(prefix, r.URL.RawPath)
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
