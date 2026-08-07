package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func TestParseKeys(t *testing.T) {
	ks := ParseKeys(" k1 , ,k2,k1 ")
	if len(ks) != 2 {
		t.Fatalf("expected 2 keys, got %d (%v)", len(ks), ks)
	}
	if !ks.has("k1") || !ks.has("k2") {
		t.Fatalf("expected k1 and k2 present: %v", ks)
	}
	if ks.has("") || ks.has("k3") {
		t.Fatalf("unexpected key membership")
	}
}

func TestRequireAuth(t *testing.T) {
	keys := ParseKeys("secret-a,secret-b")

	cases := []struct {
		name       string
		keys       KeySet
		authHeader string
		wantStatus int
	}{
		{"valid key A", keys, "Bearer secret-a", http.StatusOK},
		{"valid key B", keys, "Bearer secret-b", http.StatusOK},
		{"case-insensitive scheme", keys, "bearer secret-a", http.StatusOK},
		{"wrong key", keys, "Bearer nope", http.StatusUnauthorized},
		{"no header", keys, "", http.StatusUnauthorized},
		{"empty token", keys, "Bearer ", http.StatusUnauthorized},
		{"wrong scheme", keys, "Basic secret-a", http.StatusUnauthorized},
		{"raw key without scheme", keys, "secret-a", http.StatusUnauthorized},
		{"fail closed when unconfigured", KeySet{}, "Bearer secret-a", http.StatusUnauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := RequireAuth(tc.keys, okHandler())
			req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantStatus == http.StatusUnauthorized {
				if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", ct)
				}
				if wa := rec.Header().Get("WWW-Authenticate"); wa != "Bearer" {
					t.Errorf("WWW-Authenticate = %q, want Bearer", wa)
				}
			}
		})
	}
}

func TestMuxHealthIsPublic(t *testing.T) {
	mux := NewMux(KeySet{})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want 200", rec.Code)
	}
}

func TestMuxChatCompletionsRequiresAuth(t *testing.T) {
	mux := NewMux(ParseKeys("secret"))

	// no auth → 401
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d, want 401", rec.Code)
	}

	// valid auth → passes middleware (stub returns 501 Not Implemented)
	req = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("authenticated status = %d, want 501", rec.Code)
	}
}
