# Spec — ModelHarbor

## プロジェクトの目的

要求内容・コスト・速度・リスクに応じて AI モデルを選択するゲートウェイ。
OpenAI 互換 API で受け取り、ポリシーエンジンがモデルを選択・フォールバック・コスト記録する。

## 利用イメージ

```bash
# OpenAI 互換エンドポイント
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer $KEY" \
  -d '{"model":"auto","messages":[{"role":"user","content":"Explain SQL joins"}]}'
# → ポリシーが "SQL に強いモデル" を自動選択
```

## ルーティングポリシー例

```yaml
policies:
  - name: high-risk
    condition: contains(prompt, "個人情報")
    action: multi-vote   # 複数モデルで合議
  - name: long-context
    condition: token_count > 50000
    action: route_to: claude-3-5-sonnet
  - name: low-cost
    condition: task_type == "summarize"
    action: route_to: haiku
  - name: fallback
    action: route_to: gpt-4o-mini
```

## MVP の境界線

### やること (Phase 1〜4)

- OpenAI 互換 `/v1/chat/completions` エンドポイント
- プロバイダー 2 種類 (OpenAI / Anthropic)
- シンプルなルーティングルール (YAML)
- リトライ + フォールバック
- コスト記録 (tokens × unit price)
- ストリーミング応答 (SSE)

### やらないこと (Phase 1)

- Semantic Cache
- A/B テスト
- 個人情報マスキング
- MCP 連携

## 成功条件

| Phase   | 完成条件                              |
| ------- | ------------------------------------- |
| Phase 1 | OpenAI / Anthropic へのプロキシが動く |
| Phase 2 | YAML ルールに基づくモデル選択が動く   |
| Phase 3 | リトライ + フォールバックが動く       |
| Phase 4 | コスト記録を PostgreSQL に保存        |
| Phase 5 | SSE ストリーミング応答                |
| Phase 6 | Redis による Semantic Cache           |
