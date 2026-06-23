# Implementation Guide — ModelHarbor

## Phase 1: OpenAI / Anthropic プロキシ（1週）

- `internal/provider/openai.go` — OpenAI `/chat/completions` へ転送
- `internal/provider/anthropic.go` — Anthropic Messages API へ転送
- レスポンスを OpenAI 互換形式に正規化

**完成条件**: `curl` で OpenAI / Anthropic 両方にリクエストが通る

---

## Phase 2: YAML ルール → モデル選択（1週）

- `internal/policy/loader.go` — `policies.yaml` を読み込む
- `internal/router/router.go` — `RoutingContext` を評価し `RouteDecision` を返す
- タスク分類 (keyword-based): "SELECT" → code, "要約" → summarize

**完成条件**: `model: auto` で適切なモデルが選ばれる

---

## Phase 3: リトライ + フォールバック（3日）

- プロバイダーエラー時に `Fallbacks` リストの次のモデルへ切り替え
- 3 回失敗で `503` を返す
- リトライ間隔: exponential backoff

---

## Phase 4: コスト記録（3日）

- `internal/recorder/cost.go` — tokens × unit_price を PostgreSQL に INSERT
- 単価テーブルを `config.yaml` で管理

---

## Phase 5: SSE ストリーミング（1週）

- プロバイダーから chunked response を受け取り SSE で転送
- `Transfer-Encoding: chunked` + `Content-Type: text/event-stream`

---

## Phase 6: Redis Semantic Cache（1〜2週）

- 埋め込みモデルでプロンプトをベクトル化
- Redis Vector Search で類似キャッシュを検索 (cosine ≥ 0.95 でヒット)
