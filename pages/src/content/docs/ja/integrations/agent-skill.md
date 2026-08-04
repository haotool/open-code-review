---
title: Agent Skill
sidebar:
  order: 1
---

> **Delegate Edition（本 Fork）：** ソースから **`ocr-delegate`** をビルドした後、**`make install-skill`** で Skill をインストールします。標準ワークフローは [デリゲーションモード](../delegate/) を参照。

## Delegate Edition — Agent Skill

本 Fork は
[`skills/open-code-review/SKILL.md`](https://github.com/haotool/open-code-review-delegate/blob/main/skills/open-code-review/SKILL.md)
にデリゲート優先の SKILL を同梱しています。**`ocr-delegate`** がファイル選定とルール解決を担当し、ホストエージェントがサブスクリプション LLM でレビューを実行します。

### 前提条件

- **Git ≥ 2.41**
- サブスクリプション LLM 付き AI コーディングエージェント（Cursor、Claude Code、Codex 等）

### インストール

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

確認：

```bash
ocr-delegate -h
which ocr-delegate
```

Claude Code ユーザーは Skill ディレクトリを指定できます：

```bash
make install-skill SKILL_DIR=~/.claude/skills
```

### ワークフロー

Skill は 5 ステップのループを駆動 — エンジニアリングは `ocr-delegate`、推論はホストエージェント：

1. **Preview** — `ocr-delegate preview` でレビュー対象ファイルと mode/ref メタデータを取得。
2. **Rules** — `ocr-delegate rule <paths…>` でルールを内容別にグループ化。
3. **Diffs** — preview の mode/ref に基づき git で diff を取得（[デリゲーションモード](../delegate/) 参照）。
4. **Review** — ホストエージェントがサブスクリプション LLM で各ファイルをレビュー。
5. **Report** — Skill スキーマに従い重大度（Critical/High/Medium/Low）を分類。

例（任意の Git リポジトリ）：

```bash
cd path/to/your-repo
ocr-delegate preview -b "PR コンテキスト"
ocr-delegate rule internal/api/handler.go
```

フラグ一覧とセキュリティ規律（T7/T9）は Skill マニフェストと [デリゲーションモード](../delegate/) を参照。

### プラグインラッパー

`plugins/open-code-review/` 配下の OpenCode、Claude Code、Codex プラグインエントリは同一 Skill を指します。Slash コマンドは [Claude Code プラグイン](../claude-code/) を参照。

## 関連

- [デリゲーションモード](../delegate/) — 本 Fork の標準ワークフロー。
- [クイックスタート](../../quickstart/) — ビルド、install-skill、初回 preview。
- [インストール](../../installation/) — 前提条件と Skill インストールパス。

---

## 上流 Agent Skill（本 Fork では未提供）

> 以下は **上流 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)** のみに適用されます。

上流は完全な `ocr` CLI を呼び出す自己完結型 Skill（初回実行時の npm ベース CLI 自動インストールを含む）を提供しています。`npx skills add`、Anthropic Agent SDK 統合、内蔵 LLM レビューワークフローは上流リポジトリのドキュメントを参照してください。
