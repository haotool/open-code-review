---
title: CI/CD
sidebar:
  order: 4
---

> **Delegate Edition（本 Fork）：** 上流 `ocr` CLI を npm インストールする CI パイプラインは本リポジトリの**標準ワークフローではありません**。ソースから **`ocr-delegate`** をビルドし、Agent Skill でローカルレビュー — [デリゲーションモード](../delegate/) を参照。

## Delegate Edition — CI/CD

本 Fork は **npm ベースの CI レシピを同梱しません**。`ocr-delegate` は決定論的なファイル選定とルール解決のためのゼロネットワークバイナリです。LLM 推論はホストエージェントがサブスクリプション枠で実行します。

### 推奨アプローチ

開発フローまたはエージェントセッションでコードレビューを実行：

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

レビュー対象リポジトリで：

```bash
ocr-delegate preview --from origin/main --to HEAD
ocr-delegate rule <preview-出力のパス>
```

preview の mode/ref に基づき git で diff を取得。ホストエージェント（Cursor、Claude Code、Codex 等）がレビューを完了 — [Skill ワークフロー](../delegate/) を参照。

### オプション：ソースビルド付きセルフホスト Runner

npm なしで自動化が必要な場合、Runner 上で `ocr-delegate` をビルドし、ジョブスクリプトで `preview` / `rule` を呼び出します。**`ocr-delegate` 単体で JSON レビューコメントは生成されません** — 後処理と LLM レビューはホストエージェントの責務です。

GitHub Actions の例：

```yaml
- name: Build ocr-delegate
  run: make build
- name: Preview changed files
  run: ./dist/ocr-delegate preview --from origin/${{ github.base_ref }} --to HEAD
```

必要に応じてエージェントパイプラインを拡張してください。

## 関連

- [デリゲーションモード](../delegate/) — preview → rule → レビューワークフロー。
- [クイックスタート](../../quickstart/) — 3 ステップインストール。
- [E2E dry-run](https://github.com/haotool/open-code-review-delegate/tree/main/test/e2e) — ゼロアウトバウンド検証。

---

## 上流 CI/CD（本 Fork では未提供）

> 以下は **上流 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)** のみに適用されます。

上流は `ocr` CLI のインストール、Secrets からの LLM 設定、`ocr review --format json` の実行、PR/MR へのインラインコメント投稿を行う GitHub Actions と GitLab CI パイプラインを提供しています。

- GitHub Actions：[`examples/github_actions/ocr-review.yml`](https://github.com/alibaba/open-code-review/blob/main/examples/github_actions/ocr-review.yml)
- GitLab CI：[`examples/gitlab_ci/.gitlab-ci.yml`](https://github.com/alibaba/open-code-review/blob/main/examples/gitlab_ci/.gitlab-ci.yml)

Secrets、カスタマイズ、トラブルシューティングは上流ドキュメントを参照してください。
