# model-harbor

Multi-model AI routing gateway with OPA-based policy engine (Go).

Routes LLM requests to the optimal provider (Anthropic Claude / OpenAI-compatible) based on task type, token count, and risk level — enforced by Open Policy Agent rules.

## Architecture

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

## Usage

```bash
go build ./cmd/modelharbor

# Route a request
curl -X POST http://localhost:8080/v1/chat \
  -H "Authorization: Bearer <key>" \
  -H "Content-Type: application/json" \
  -d '{"task_type":"coding","messages":[{"role":"user","content":"Write a binary search in Go"}]}'
```

## Tech Stack

- **Language:** Go 1.22
- **Policy engine:** Open Policy Agent (OPA)
- **Providers:** Anthropic Claude, OpenAI-compatible APIs
- **Auth:** Bearer token (constant-time comparison via `crypto/subtle`)
- **Cost tracking:** token × unit_price → PostgreSQL

## API

| Method | Path | Description |
|---|---|---|
| POST | `/v1/chat` | Route a chat request to the optimal model |
| GET | `/v1/health` | Health check |

## Architecture Decision Records

- [ADR-001](docs/adr/ADR-001-openai-compatible.md): OpenAI-compatible API surface
- [ADR-002](docs/adr/ADR-002-opa-policy.md): OPA for routing policy

## Status

Core routing and auth complete. OPA policy integration and cost recording in progress.

## License

MIT

