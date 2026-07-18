package provider_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipslidersand/model-harbor/internal/provider"
)

func TestAnthropicProvider_Complete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "ant-key" {
			t.Errorf("missing anthropic api key header")
		}
		// Return Anthropic-format response
		resp := map[string]any{
			"id":          "msg_test",
			"model":       "claude-3-5-haiku-20241022",
			"stop_reason": "end_turn",
			"content":     []map[string]any{{"type": "text", "text": "こんにちは！"}},
			"usage":       map[string]any{"input_tokens": 8, "output_tokens": 4},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := provider.NewAnthropic("ant-key", srv.URL)
	req := &provider.ChatRequest{
		Model:    "claude-3-5-haiku-20241022",
		Messages: []provider.Message{{Role: "user", Content: "こんにちは"}},
	}

	got, err := p.Complete(req)
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if got.Object != "chat.completion" {
		t.Errorf("Object: got %s, want chat.completion", got.Object)
	}
	if len(got.Choices) == 0 || got.Choices[0].Message.Content != "こんにちは！" {
		t.Errorf("unexpected content: %+v", got.Choices)
	}
	if got.Usage.PromptTokens != 8 {
		t.Errorf("PromptTokens: got %d, want 8", got.Usage.PromptTokens)
	}
}

func TestAnthropicProvider_SystemMessageExtracted(t *testing.T) {
	var capturedBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedBody)
		resp := map[string]any{
			"id": "msg_sys", "model": "claude-3-5-haiku-20241022", "stop_reason": "end_turn",
			"content": []map[string]any{{"type": "text", "text": "ok"}},
			"usage":   map[string]any{"input_tokens": 1, "output_tokens": 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := provider.NewAnthropic("k", srv.URL)
	req := &provider.ChatRequest{
		Model: "claude-3-5-haiku-20241022",
		Messages: []provider.Message{
			{Role: "system", Content: "You are helpful."},
			{Role: "user", Content: "Hi"},
		},
	}
	if _, err := p.Complete(req); err != nil {
		t.Fatal(err)
	}

	if capturedBody["system"] != "You are helpful." {
		t.Errorf("system field: got %v", capturedBody["system"])
	}
	msgs, ok := capturedBody["messages"].([]any)
	if !ok || len(msgs) != 1 {
		t.Errorf("expected 1 message (system excluded), got %v", capturedBody["messages"])
	}
}
