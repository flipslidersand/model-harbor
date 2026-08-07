// Package httpapi provides the OpenAI-compatible HTTP surface for ModelHarbor
// and the middleware that guards it.
package httpapi

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// envAPIKeys is the environment variable holding the comma-separated list of
// accepted API keys (Bearer tokens) for the gateway.
const envAPIKeys = "MODELHARBOR_API_KEYS"

// KeySet is the set of accepted API keys. An empty set means the gateway is
// unconfigured and — by design — rejects every request (fail closed) so a
// public deployment never proxies to paid backends without a caller key
// (Issue #3).
type KeySet map[string]struct{}

// LoadKeysFromEnv builds a KeySet from MODELHARBOR_API_KEYS. Keys are
// comma-separated; surrounding whitespace and empty entries are ignored.
func LoadKeysFromEnv() KeySet {
	return ParseKeys(os.Getenv(envAPIKeys))
}

// ParseKeys parses a comma-separated key list into a KeySet.
func ParseKeys(raw string) KeySet {
	keys := KeySet{}
	for _, part := range strings.Split(raw, ",") {
		k := strings.TrimSpace(part)
		if k != "" {
			keys[k] = struct{}{}
		}
	}
	return keys
}

// has reports whether token matches any configured key using a constant-time
// comparison to avoid leaking key material through timing.
func (ks KeySet) has(token string) bool {
	var ok bool
	for k := range ks {
		if subtle.ConstantTimeCompare([]byte(k), []byte(token)) == 1 {
			ok = true
		}
	}
	return ok
}

// bearerToken extracts the token from an `Authorization: Bearer <token>`
// header. The scheme match is case-insensitive per RFC 7235. It returns the
// token and whether a well-formed Bearer credential was present.
func bearerToken(h http.Header) (string, bool) {
	auth := h.Get("Authorization")
	if auth == "" {
		return "", false
	}
	const prefix = "bearer "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(auth[len(prefix):])
	if token == "" {
		return "", false
	}
	return token, true
}

// RequireAuth wraps next so that only requests carrying a valid
// `Authorization: Bearer <key>` header (matching keys) are served. All other
// requests receive 401 with an OpenAI-style error body. When keys is empty the
// middleware fails closed and rejects everything.
func RequireAuth(keys KeySet, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header)
		if !ok {
			writeAuthError(w, http.StatusUnauthorized,
				"Missing or malformed Authorization header. Expected 'Authorization: Bearer <API key>'.",
				"missing_api_key")
			return
		}
		if !keys.has(token) {
			writeAuthError(w, http.StatusUnauthorized,
				"Incorrect API key provided.", "invalid_api_key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// writeAuthError emits an OpenAI-compatible error envelope.
func writeAuthError(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"message": message,
			"type":    "invalid_request_error",
			"code":    code,
			"param":   nil,
		},
	})
}
