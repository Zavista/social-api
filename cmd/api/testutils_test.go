package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zavista/social-api/internal/auth"
	"github.com/zavista/social-api/internal/store"
	"github.com/zavista/social-api/internal/store/cache"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	return &application{
		logger:        slog.New(slog.DiscardHandler),
		store:         store.NewMockStore(),
		cacheStorage:  cache.NewMockStore(),
		authenticator: &auth.TestAuthenticator{},
		config:        cfg,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d, got %d", expected, actual)
	}
}
