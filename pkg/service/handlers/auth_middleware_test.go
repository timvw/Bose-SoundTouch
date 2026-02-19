package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthMiddleware_ValidCredentials(t *testing.T) {
	handler := BasicAuthMiddleware("admin", "secret123")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/mgmt/test", nil)
	req.SetBasicAuth("admin", "secret123")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got %q", rr.Body.String())
	}
}

func TestBasicAuthMiddleware_WrongUsername(t *testing.T) {
	handler := BasicAuthMiddleware("admin", "secret123")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with wrong username")
	}))

	req := httptest.NewRequest(http.MethodGet, "/mgmt/test", nil)
	req.SetBasicAuth("wrong", "secret123")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if rr.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header to be set")
	}
}

func TestBasicAuthMiddleware_WrongPassword(t *testing.T) {
	handler := BasicAuthMiddleware("admin", "secret123")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with wrong password")
	}))

	req := httptest.NewRequest(http.MethodGet, "/mgmt/test", nil)
	req.SetBasicAuth("admin", "wrongpass")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if rr.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header to be set")
	}
}

func TestBasicAuthMiddleware_MissingAuthHeader(t *testing.T) {
	handler := BasicAuthMiddleware("admin", "secret123")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called without auth header")
	}))

	req := httptest.NewRequest(http.MethodGet, "/mgmt/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if rr.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header to be set")
	}
}

func TestBasicAuthMiddleware_EmptyCredentials(t *testing.T) {
	handler := BasicAuthMiddleware("admin", "secret123")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called with empty credentials")
	}))

	req := httptest.NewRequest(http.MethodGet, "/mgmt/test", nil)
	req.SetBasicAuth("", "")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if rr.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header to be set")
	}
}
