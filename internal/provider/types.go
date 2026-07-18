package provider

// ChatRequest は OpenAI 互換リクエスト形式。
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Stream    bool      `json:"stream"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse は OpenAI 互換レスポンス形式。
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Provider はバックエンド LLM プロバイダーのインターフェース。
type Provider interface {
	// Name はプロバイダー名を返す ("openai" | "anthropic")。
	Name() string
	// Complete は ChatRequest を送信し OpenAI 互換レスポンスを返す。
	Complete(req *ChatRequest) (*ChatResponse, error)
}
