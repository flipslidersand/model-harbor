package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipslidersand/model-harbor/internal/provider"
	"github.com/flipslidersand/model-harbor/internal/router"
)

// mockProvider は Provider インターフェースのスタブ。
type mockProvider struct {
	name string
	resp *provider.ChatResponse
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) Complete(req *provider.ChatRequest) (*provider.ChatResponse, error) {
	resp := *m.resp
	resp.Model = req.Model
	return &resp, nil
}

var mockResp = &provider.ChatResponse{
	ID: "mock", Object: "chat.completion",
	Choices: []provider.Choice{{Message: provider.Message{Role: "assistant", Content: "ok"}}},
}

func TestRouter_RoutesGPTToOpenAI(t *testing.T) {
	oai := &mockProvider{name: "openai", resp: mockResp}
	ant := &mockProvider{name: "anthropic", resp: mockResp}
	rt := router.New(oai, ant)

	body, _ := json.Marshal(provider.ChatRequest{Model: "gpt-4o-mini", Messages: []provider.Message{{Role: "user", Content: "hi"}}})
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	rt.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp provider.ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Model != "gpt-4o-mini" {
		t.Errorf("expected model gpt-4o-mini, got %s", resp.Model)
	}
}

func TestRouter_RoutesClaudeToAnthropic(t *testing.T) {
	var calledProvider string
	oai := &mockProvider{name: "openai", resp: mockResp}
	ant := &mockProvider{name: "anthropic", resp: mockResp}
	_ = calledProvider

	rt := router.New(oai, ant)

	body, _ := json.Marshal(provider.ChatRequest{Model: "claude-3-5-haiku-20241022", Messages: []provider.Message{{Role: "user", Content: "hi"}}})
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	rt.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var resp provider.ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Model != "claude-3-5-haiku-20241022" {
		t.Errorf("expected claude model, got %s", resp.Model)
	}
	_ = ant
}

func TestRouter_404OnUnknownPath(t *testing.T) {
	rt := router.New(&mockProvider{name: "openai", resp: mockResp}, &mockProvider{name: "anthropic", resp: mockResp})
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	rt.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestRouter_405OnGet(t *testing.T) {
	rt := router.New(&mockProvider{name: "openai", resp: mockResp}, &mockProvider{name: "anthropic", resp: mockResp})
	req := httptest.NewRequest(http.MethodGet, "/v1/chat/completions", nil)
	w := httptest.NewRecorder()
	rt.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", w.Code)
	}
}
