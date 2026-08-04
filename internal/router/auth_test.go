package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipslidersand/model-harbor/internal/router"
)

func TestBearerAuth_AllowsWhenKeyEmpty(t *testing.T) {
	called := false
	h := router.BearerAuth("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !called {
		t.Fatal("handler should be called when apiKey is empty")
	}
}

func TestBearerAuth_RejectsNoHeader(t *testing.T) {
	h := router.BearerAuth("secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

func TestBearerAuth_RejectsWrongToken(t *testing.T) {
	h := router.BearerAuth("correct", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	h.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

func TestBearerAuth_AllowsCorrectToken(t *testing.T) {
	called := false
	h := router.BearerAuth("mysecret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer mysecret")
	h.ServeHTTP(w, req)
	if !called {
		t.Fatal("handler should be called with correct token")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
}
