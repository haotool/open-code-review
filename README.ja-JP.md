<div align="center">
  <h1>Open Code Review — Delegate Edition</h1>
  <p><strong><a href="https://github.com/haotool">haotool</a> によるセキュリティ強化 Fork</strong></p>
</div>

<p align="center">
  <a href="README.md">English</a> | <a href="README.zh-CN.md">简体中文</a> | <a href="README.zh-TW.md">繁體中文</a> | 日本語 | <a href="README.ko-KR.md">한국어</a> | <a href="README.ru-RU.md">Русский</a>
</p>

<p align="center">
  <a href="https://github.com/haotool/open-code-review-delegate/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/haotool/open-code-review-delegate/ci.yml?style=flat-square" /></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/haotool/open-code-review-delegate?style=flat-square" /></a>
  <a href="https://github.com/alibaba/open-code-review"><img alt="Upstream" src="https://img.shields.io/badge/upstream-alibaba%2Fopen--code--review-blue?style=flat-square" /></a>
</p>

---

## 重要なお知らせ

本リポジトリは Alibaba 公式製品**ではありません**。[haotool](https://github.com/haotool) が Apache License 2.0 の下で独立維持する Fork で、[alibaba/open-code-review](https://github.com/alibaba/open-code-review) から派生しています。

- リポジトリ：[`haotool/open-code-review-delegate`](https://github.com/haotool/open-code-review-delegate)（旧 `haotool/open-code-review` URL は自動リダイレクトされます）。
- Alibaba の商標は使用せず、公式の後押しを暗示しません。
- upstream の著作権表示は [NOTICE](NOTICE) を参照してください。
- ライセンス全文は [LICENSE](LICENSE) を参照してください。

## 概要

Open Code Review (Delegate Edition) は単一 CLI バイナリ **`ocr-delegate`** を提供し、決定論的なコードレビュー工学を担当します。

- **ファイル選定** — レビュー対象の変更ファイルを正確に決定
- **ルール解決** — `rule.json` からファイルごとのレビュールールをマッチ
- **ゼロネットワーク** — バイナリは LLM 呼び出しや外向き接続を行わない
- **ゼロ資格情報** — API キー、トークン、プロバイダー設定不要

AI コーディングエージェント（Cursor、Claude Code、Codex 等）がサブスクリプション LLM で実際のレビューを実行します。`ocr-delegate` が工学、Host Agent が知能を担当します。

upstream 配布モデルの npm 自動更新、プリビルトバイナリのダウンロード、組み込みプロバイダーエンドポイントによるサプライチェーンリスクを排除します。

## クイックインストール（3 ステップ）

ソースからビルド — npm、プリビルト、API キー不要：

```bash
# 1. Clone
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate

# 2. ocr-delegate をビルド
make build

# 3. エージェント skill をインストール
make install-skill
# Claude Code ユーザー:
make install-skill SKILL_DIR=~/.claude/skills
```

検証:

```bash
which ocr-delegate && ocr-delegate -h
```

インストールは `test/e2e/dryrun.sh` で検証されます（[E2E 検証](test/e2e/README.md)）。

## 使い方

`skills/open-code-review/SKILL.md` がデリゲートワークフローの単一の正（SSOT）です。

### 1. Preview — スコープ決定

```bash
ocr-delegate preview [--from main --to feature] [--commit <hash>] [-b "context"]
```

### 2. ファイルのルール取得

```bash
ocr-delegate rule <path1> <path2> ...
```

### 3. diff 取得とレビュー

preview 出力の mode/ref メタデータに従い git を使用。Host Agent が構造化 findings を生成します。

完全なワークフロー、セキュリティ規律（T7/T9）、出力スキーマは [skill ドキュメント](skills/open-code-review/SKILL.md) を参照してください。

Codex、Claude Code、Cursor の正式サポートインストールパスは [エージェントサポートマトリクス](docs/AGENT_SUPPORT.md) を参照してください。

## セキュリティ posture

| 属性 | 保証 |
|------|------|
| ゼロ外向き | ネットワーク接続なし（CI 文字列スキャン + E2E で検証） |
| ゼロ資格情報 | Delegate モードで LLM プロバイダー設定不要 |
| 最小依存クロージャ | 禁止モジュール（`internal/llm`、`telemetry`、`mcp`、`viewer`、`session`）を除外 |
| ソース優先配布 | npm ラッパー、プリビルトダウンロード、自動更新なし |
| サプライチェーン除去 | upstream インストールスクリプト、npm パッケージ、GitHub Action を削除 |

詳細: [SECURITY.md](SECURITY.md)

## ビルド再現と検証

本 Fork は**ソース優先**で、プリビルトバイナリは配布しません。

```bash
make build                          # dist/ocr-delegate を生成
shasum -a 256 dist/ocr-delegate     # ローカルチェックサム

# クロスプラットフォームローカルビルド（任意）
make dist                           # 全プラットフォーム + sha256sum.txt
```

## 開発

```bash
make build       # ocr-delegate をビルド
make test        # テストスイート実行
make check       # fmt + vet + mod tidy
make coverage    # カバレッジレポート（80% 閾値）
```

[CONTRIBUTING.ja-JP.md](CONTRIBUTING.ja-JP.md) を参照してください。

## 上流レガシーディレクトリ

upstream から残存するが、本 Fork では**メンテナンスしない**ディレクトリ:

| ディレクトリ | 状態 |
|------|------|
| `pages/` | upstream ランディング — 本 Fork ではビルド/デプロイしない |
| `extensions/` | upstream VS Code 拡張 — `ocr-delegate` ではなく完全な `ocr` CLI が必要 |

## ライセンスと帰属

[Apache-2.0](LICENSE)。upstream Alibaba 貢献者の著作権表示は [NOTICE](NOTICE) に保持されています。

本 Fork は upstream を大幅に変更しています。v1.0.0 の変更は [CHANGELOG.md](CHANGELOG.md) を参照してください。

## リリースチェックリスト (S1–S7)

| ID | 要件 | 検証 |
|----|------|------|
| S1 | ゼロネットワーク | `test/e2e/dryrun.sh` + CI Gate 4 文字列スキャン |
| S2 | ゼロ資格情報 | E2E dry-run + delegate 専用バイナリ |
| S3 | 決定論エンジン | 26 シナリオ golden 等価スイート |
| S4 | ソース優先 + delegate-only skill | インストール経路に npm/postinstall サプライチェーンなし；Skill は `ocr-delegate` のみ参照（完全な `ocr` CLI / 外部 LLM API 不可） |
| S5 | インストール ≤3 ステップ | E2E dry-run 3 ステップ検証 |
| S6 | テスト通過 | CI Gate 2 (`go test ./...`) |
| S7 | 注入耐性 | 敵対 fixture ライブラリ + E2E spot-check |
