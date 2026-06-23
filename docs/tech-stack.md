# Tech Stack — ModelHarbor

## 言語・バージョン

- Go 1.22+

## 主要パッケージ

| パッケージ                         | 役割                         | 選定理由                    |
| ---------------------------------- | ---------------------------- | --------------------------- |
| `net/http` + `encoding/json`       | OpenAI 互換 API サーバー     | 標準ライブラリで十分        |
| `github.com/open-policy-agent/opa` | ポリシーエンジン             | Rego でルールを宣言的に定義 |
| `github.com/redis/go-redis/v9`     | Semantic Cache (Phase 6)     | Redis クライアント標準      |
| `github.com/jackc/pgx/v5`          | PostgreSQL コスト記録        | 高速・context 対応          |
| `go.opentelemetry.io/otel`         | コスト・レイテンシメトリクス | 既存監視基盤に統合          |
| `github.com/spf13/cobra`           | CLI                          | Go 標準 CLI フレームワーク  |
| `go.uber.org/zap`                  | 構造化ログ                   | リクエストトレーサビリティ  |

## アーキテクチャ

```
Client (OpenAI SDK)
  ↓ POST /v1/chat/completions
[Gateway Handler]
  ↓
[Task Classifier]   prompt → task_type (summarize / code / qa / ...)
  ↓
[Policy Engine (OPA)]  task_type + token_count + risk → model 選択
  ↓
[Provider Router]
  ├── OpenAI Provider
  └── Anthropic Provider
  ↓ (SSE stream)
[Cost Recorder]     tokens × unit_price → PostgreSQL
  ↓
Client
```
