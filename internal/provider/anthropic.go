package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AnthropicProvider は Anthropic Messages API へのプロキシ。
// OpenAI 互換形式 → Anthropic 形式に変換して転送し、レスポンスを正規化する。
type AnthropicProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewAnthropic(apiKey, baseURL string) *AnthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	return &AnthropicProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (p *AnthropicProvider) Name() string { return "anthropic" }

// anthropicRequest は Anthropic Messages API のリクエスト形式。
type anthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []anthropicMessage `json:"messages"`
	System    string             `json:"system,omitempty"`
	MaxTokens int                `json:"max_tokens"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	ID      string              `json:"id"`
	Model   string              `json:"model"`
	Content []anthropicContent  `json:"content"`
	Usage   anthropicUsage      `json:"usage"`
	StopReason string           `json:"stop_reason"`
}

type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func (p *AnthropicProvider) Complete(req *ChatRequest) (*ChatResponse, error) {
	aReq := p.toAnthropicRequest(req)
	body, err := json.Marshal(aReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic: status %d: %s", resp.StatusCode, respBody)
	}

	var aResp anthropicResponse
	if err := json.Unmarshal(respBody, &aResp); err != nil {
		return nil, err
	}
	return p.toOpenAIResponse(&aResp), nil
}

// toAnthropicRequest は OpenAI 形式を Anthropic Messages 形式に変換する。
func (p *AnthropicProvider) toAnthropicRequest(req *ChatRequest) *anthropicRequest {
	aReq := &anthropicRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
	}
	if aReq.MaxTokens == 0 {
		aReq.MaxTokens = 1024
	}

	for _, m := range req.Messages {
		if m.Role == "system" {
			aReq.System = m.Content
		} else {
			aReq.Messages = append(aReq.Messages, anthropicMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	}
	return aReq
}

// toOpenAIResponse は Anthropic レスポンスを OpenAI 互換形式に正規化する。
func (p *AnthropicProvider) toOpenAIResponse(aResp *anthropicResponse) *ChatResponse {
	text := ""
	if len(aResp.Content) > 0 {
		text = aResp.Content[0].Text
	}
	finishReason := "stop"
	if aResp.StopReason == "max_tokens" {
		finishReason = "length"
	}
	return &ChatResponse{
		ID:     aResp.ID,
		Object: "chat.completion",
		Model:  aResp.Model,
		Choices: []Choice{
			{
				Index:        0,
				Message:      Message{Role: "assistant", Content: text},
				FinishReason: finishReason,
			},
		},
		Usage: Usage{
			PromptTokens:     aResp.Usage.InputTokens,
			CompletionTokens: aResp.Usage.OutputTokens,
			TotalTokens:      aResp.Usage.InputTokens + aResp.Usage.OutputTokens,
		},
	}
}
