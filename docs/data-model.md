# Data Model — ModelHarbor

## リクエスト・レスポンス

```go
// OpenAI 互換リクエスト
type ChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
    Stream   bool      `json:"stream"`
    MaxTokens int      `json:"max_tokens,omitempty"`
}

type Message struct {
    Role    string `json:"role"`    // "user" | "assistant" | "system"
    Content string `json:"content"`
}

// ポリシー判定用コンテキスト
type RoutingContext struct {
    Request    ChatRequest
    TaskType   string  // "code" | "summarize" | "qa" | "translation"
    TokenCount int
    RiskScore  float64 // 0.0〜1.0 (個人情報・機密語を検出)
}

// 選択されたプロバイダー
type RouteDecision struct {
    Provider  string // "openai" | "anthropic"
    Model     string // "gpt-4o" | "claude-3-5-sonnet-20241022"
    Reason    string
    Fallbacks []RouteDecision
}
```

## コスト記録 (PostgreSQL)

```sql
CREATE TABLE requests (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    model        TEXT NOT NULL,
    provider     TEXT NOT NULL,
    input_tokens  INTEGER NOT NULL,
    output_tokens INTEGER NOT NULL,
    cost_usd     NUMERIC(10, 6) NOT NULL,
    latency_ms   INTEGER NOT NULL,
    task_type    TEXT
);
```

## ルーティングポリシー (YAML → OPA Rego)

```yaml
# policies.yaml → Rego に変換
policies:
  - name: high-risk
    condition: risk_score > 0.8
    action: multi_vote
    models: [gpt-4o, claude-3-5-sonnet]
  - name: long-context
    condition: token_count > 50000
    action: route_to
    model: claude-3-5-sonnet-20241022
  - name: default
    action: route_to
    model: gpt-4o-mini
```
