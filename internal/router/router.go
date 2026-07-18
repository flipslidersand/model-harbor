package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flipslidersand/model-harbor/internal/provider"
)

// Router はモデル名でプロバイダーを選択し /v1/chat/completions を処理する。
type Router struct {
	openai    provider.Provider
	anthropic provider.Provider
}

func New(openai, anthropic provider.Provider) *Router {
	return &Router{openai: openai, anthropic: anthropic}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/v1/chat/completions" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var chatReq provider.ChatRequest
	if err := json.NewDecoder(req.Body).Decode(&chatReq); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	p := r.selectProvider(chatReq.Model)
	resp, err := p.Complete(&chatReq)
	if err != nil {
		http.Error(w, "provider error: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (r *Router) selectProvider(model string) provider.Provider {
	if strings.HasPrefix(model, "claude") {
		return r.anthropic
	}
	return r.openai
}
