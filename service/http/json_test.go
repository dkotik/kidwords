package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/kidwords/service"
	"github.com/dkotik/kidwords/service/secret/mock"
)

func TestHandlersJSON(t *testing.T) {
	const prefix = "/"
	repository := mock.New()
	service, err := service.New(
		mockAuthenticator{},
		repository,
	)
	if err != nil {
		t.Fatal(err)
	}
	mux, err := NewJSON(
		service,
		WithPathPrefix(prefix),
		WithAdaptor(
			htadaptor.New(
				htadaptor.WithErrorHandler(htadaptor.ErrorHandlerFunc(
					func(w http.ResponseWriter, r *http.Request, err error) error {
						t.Log("request:", r.Method, r.URL.String())
						t.Fatal(err)
						return nil
					})),
			),
		),
	)
	if err != nil {
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

	const testKeyName = "test-key"
	var testKeyID string
	t.Run("createPaperKey", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", testKeyName)
		req, err := http.NewRequest("POST", prefix, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		data, sc, err := server(req)
		if err != nil {
			t.Fatal(err)
		}
		if sc != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, sc)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty body")
		}
		keys, err := repository.List(t.Context(), mockUserID)
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(keys))
		}
		if keys[0].Name != testKeyName {
			t.Fatalf("expected key name %s, got %s", testKeyName, keys[0].Name)
		}
		testKeyID = keys[0].ID
	})

	t.Run("updatePaperKey", func(t *testing.T) {
		form := url.Values{}
		form.Set("id", testKeyID)
		form.Set("name", testKeyName+":updated")
		req, err := http.NewRequest("PUT", prefix, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		data, sc, err := server(req)
		if err != nil {
			t.Fatal(err)
		}
		if sc != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, sc)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty body")
		}
	})

	t.Run("listPaperKeys", func(t *testing.T) {
		req, err := http.NewRequest("GET", prefix, nil)
		if err != nil {
			t.Fatal(err)
		}
		data, sc, err := server(req)
		if err != nil {
			t.Fatal(err)
		}
		if sc != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, sc)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty body")
		}
		keys, err := repository.List(t.Context(), mockUserID)
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(keys))
		}
		if keys[0].Name != testKeyName+":updated" {
			t.Fatalf("expected key name %s, got %s", testKeyName+":updated", keys[0].Name)
		}
	})

	t.Run("deletePaperKey", func(t *testing.T) {
		req, err := http.NewRequest(
			"DELETE",
			fmt.Sprintf("%s?delete=%s", prefix, testKeyID),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}
		data, sc, err := server(req)
		if err != nil {
			t.Fatal(err)
		}
		if sc != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, sc)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty body")
		}
		keys, err := repository.List(t.Context(), mockUserID)
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 0 {
			t.Fatalf("expected 0 keys, got %d", len(keys))
		}
	})
}
