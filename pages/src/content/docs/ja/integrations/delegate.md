---
title: デリゲーションモード
sidebar:
  order: 5
---

> **Delegate Edition（haotool fork）：** 本リポジトリの標準ワークフロー。ソースから `ocr-delegate` をビルド — npm 不要、OCR 側に LLM API キー不要。

`ocr-delegate` が確定的エンジニアリング（ファイル選択、ルール解決）を担当し、ホストエージェント（Cursor、Claude Code、Codex 等）がサブスクリプション LLM でレビューを実行します。バイナリは外部ネットワークに接続しません。

## デリゲーションモードを使う場面

サブスクリプション型 AI コーディングエージェントで、ホストエージェントに LLM が既に組み込まれている場合向けです。別途モデルエンドポイントを設定せず、ホストエージェントのクォータをレビュー推論に再利用できます。

次の場合に適しています：

1. サブスクリプション制 AI コーディングエージェントで、既存クォータをコードレビューに再利用したい（追加 API キーやモデル設定不要）。
2. OCR は工程スキャフォールド（ファイルフィルタ、ルール解決、除外ロジック）のみ必要で、LLM 推論はホストエージェントに任せたい。
3. ファイル一覧 + ルールという構造化入力が必要なカスタムエージェントパイプラインを構築している。

## 前提条件

[README](https://github.com/haotool/open-code-review-delegate#quick-install-3-steps) からソースビルド：

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
ocr-delegate -h
```

OCR 側に LLM 設定や環境変数は不要です。

## Skill のインストール

```bash
make install-skill
make install-skill SKILL_DIR=~/.claude/skills  # Claude Code
```

正本 Skill: `skills/open-code-review/SKILL.md`。`plugins/open-code-review/` のプラグインは同 Skill を指します。

## ワークフロー

### ステップ 1：Preview — レビュー対象を決定

```bash
ocr-delegate preview [--from <ref> --to <ref>] [--commit <hash>] [--exclude <patterns>]
```

出力：

- **mode** — workspace / range / commit
- **ref メタデータ** — from、to、commit、merge\_base
- **レビュー対象ファイル一覧** — パス、status、追加/削除行数
- **除外ファイル** — 除外理由付き

よく使う呼び出し：

| シナリオ | コマンド |
|----------|---------|
| ワークスペース変更 | `ocr-delegate preview` |
| ブランチ比較 | `ocr-delegate preview --from main --to feature` |
| 単一 commit | `ocr-delegate preview -c abc123` |

### ステップ 2：ファイルのルールを取得

```bash
ocr-delegate rule <path1> <path2> ...
```

ステップ 1 のレビュー対象パスを渡します。出力はルール内容ごとにグループ化されます。

### ステップ 3：diff を取得

ステップ 1 の mode/ref 情報に基づき、git を直接使用します：

**Range モード**（merge\_base あり）：

```bash
git diff <merge_base>..<to> -- <path>
```

**Commit モード**：

```bash
git show <commit> -- <path>
```

**Workspace モード**：

```bash
git diff HEAD -- <path>        # 追跡済みファイル
cat <path>                     # 新規未追跡ファイル
```

### ステップ 4：ファイルごとにレビュー

各レビュー対象ファイルについて、ステップ 3 の diff とステップ 2 のルールグループを参照し、ホストエージェントの LLM でレビューします。

### ステップ 5：レポート

Skill スキーマに従い、重大度（Critical/High/Medium/Low）を分類します。

## サブコマンド

| コマンド | 用途 |
|---------|------|
| `ocr-delegate preview` | レビュー対象ファイルと mode/ref メタデータを一覧 |
| `ocr-delegate rule <path...>` | ルール内容ごとにグループ化して解決 |

## 共通フラグ

| フラグ | 説明 |
|--------|------|
| `--from <ref>` | Range モードのソース ref |
| `--to <ref>` | Range モードのターゲット ref |
| `-c, --commit <hash>` | 単一 commit モード |
| `--repo <path>` | リポジトリルート（デフォルト：cwd） |
| `--rule <path>` | カスタム rule.json パス |
| `--exclude <patterns>` | カンマ区切りの除外パターン |
| `-b, --background <text>` | ビジネスコンテキスト |
| `-B, --background-file <path>` | Markdown ファイルからビジネスコンテキストを読み込み |

## 上流との比較

| モード | LLM 呼び出し | 本 fork |
|--------|-------------|---------|
| Agent Skill / `ocr review` | OCR 内蔵 LLM | 非提供 — バナー付き上流ドキュメント参照 |
| **デリゲーション（`ocr-delegate`）** | ホストエージェント | **提供 — 主ワークフロー** |

## 関連

- [README](https://github.com/haotool/open-code-review-delegate)
- [E2E dry-run](https://github.com/haotool/open-code-review-delegate/tree/main/test/e2e)
