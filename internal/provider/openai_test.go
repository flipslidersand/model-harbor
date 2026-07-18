package provider_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipslidersand/model-harbor/internal/provider"
)

func TestOpenAIProvider_Complete(t *testing.T) {
	want := provider.ChatResponse{
		ID:     "chatcmpl-test",
		Object: "chat.completion",
		Model:  "gpt-4o-mini",
		Choices: []provider.Choice{
			{Index: 0, Message: provider.Message{Role: "assistant", Content: "Hello!"}, FinishReason: "stop"},
		},
		Usage: provider.Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing auth header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	p := provider.NewOpenAI("test-key", srv.URL)
	req := &provider.ChatRequest{
		Model:    "gpt-4o-mini",
		Messages: []provider.Message{{Role: "user", Content: "hi"}},
	}

	got, err := p.Complete(req)
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if got.ID != want.ID {
		t.Errorf("ID: got %s, want %s", got.ID, want.ID)
	}
	if len(got.Choices) == 0 || got.Choices[0].Message.Content != "Hello!" {
		t.Errorf("unexpected choices: %+v", got.Choices)
	}
}

func TestOpenAIProvider_UpstreamError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	p := provider.NewOpenAI("bad-key", srv.URL)
	_, err := p.Complete(&provider.ChatRequest{Model: "gpt-4o-mini"})
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}
