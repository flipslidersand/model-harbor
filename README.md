# model-harbor

Multi-model AI routing gateway with OPA-based policy engine (Go).
Routes LLM requests to the optimal provider (Anthropic Claude / OpenAI-compatible) based on task type, token count, and risk level — enforced by Open Policy Agent rules.

OPA ベースのポリシーエンジンを持つ、マルチモデル AI ルーティングゲートウェイです（Go）。
タスク種別・トークン数・リスクレベルに応じて LLM リクエストを最適なプロバイダーへルーティングします。ルールは Open Policy Agent で管理されます。

## Architecture / アーキテクチャ

```
Client → POST /v1/chat
              │
       ┌──────▼──────┐
       │   Router     │  task_type + token_count + risk → model selection
       │   (OPA)      │
       └──────┬──────┘
              │
     ┌────────┴────────┐
     ▼                 ▼
 Anthropic          OpenAI-compatible
 (Claude)           (any endpoint)
```

## Usage / 使い方

```bash
go build ./cmd/modelharbor

curl -X POST http://localhost:8080/v1/chat \
  -H "Authorization: Bearer <key>" \
  -H "Content-Type: application/json" \
  -d '{"task_type":"coding","messages":[{"role":"user","content":"Write a binary search in Go"}]}'
```

## Tech Stack / 技術スタック

- **Language / 言語:** Go 1.22
- **Policy engine / ポリシーエンジン:** Open Policy Agent (OPA)
- **Providers / プロバイダー:** Anthropic Claude, OpenAI-compatible APIs
- **Auth / 認証:** Bearer token (constant-time via `crypto/subtle`)

## API

| Method | Path | Description / 説明 |
|---|---|---|
| POST | `/v1/chat` | Route a chat request / チャットリクエストをルーティング |
| GET | `/v1/health` | Health check / ヘルスチェック |

## Status / ステータス

Core routing and auth complete. OPA policy integration in progress.
コアルーティングと認証は完成。OPA ポリシー統合は実装中。

## License

MIT
