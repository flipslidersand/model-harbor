# ADR-001: 外部 API を OpenAI 互換形式に統一する

- **日付**: 2026-06-22
- **状態**: Accepted

## 決定

クライアントには常に OpenAI 互換の `/v1/chat/completions` を公開し、
内部でプロバイダーごとのフォーマットに変換する。

## 理由

- 既存の OpenAI SDK・ツール（LangChain 等）をそのまま向けられる
- プロバイダーの追加・変更をクライアント側に影響させない
- ストリーミング (SSE) の形式も OpenAI に統一することで実装を単純化できる

## トレードオフ

- Anthropic 固有の機能（extended thinking・tool use の細かい挙動差）を完全には吸収できない
