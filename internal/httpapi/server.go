package httpapi

import (
	"encoding/json"
	"net/http"
)

// NewMux builds the ModelHarbor HTTP router. `/healthz` is public; the
// OpenAI-compatible `/v1/chat/completions` endpoint is guarded by RequireAuth
// so callers must present a valid API key (Issue #3).
//
// The chat-completions handler is a placeholder until the routing/provider
// stack lands (see docs/spec.md, Phase 1); it is wired behind auth from the
// outset so the endpoint is never exposed unauthenticated.
func NewMux(keys KeySet) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealth)
	mux.Handle("/v1/chat/completions", RequireAuth(keys, http.HandlerFunc(handleChatCompletions)))
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleChatCompletions is a stub: authentication is enforced by the
// middleware; the routing/provider logic is not yet implemented.
func handleChatCompletions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"message": "chat completions routing is not implemented yet",
			"type":    "server_error",
			"code":    "not_implemented",
		},
	})
}
