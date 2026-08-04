package router

import (
	"net/http"
	"strings"
)

// BearerAuth returns middleware that rejects requests without a valid Bearer token.
// If apiKey is empty, all requests are allowed (dev mode).
func BearerAuth(apiKey string, next http.Handler) http.Handler {
	if apiKey == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token, ok := strings.CutPrefix(auth, "Bearer ")
		if !ok || token != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
