---
title: クイックスタート
sidebar:
  order: 3
---

数分で初回のコードレビューを実行できます。

> **Delegate Edition（本 Fork）：** 標準ワークフローはソースからビルドした **`ocr-delegate`** を使用 — npm 不要、OCR 側に LLM API キー不要。完全な手順は [デリゲーションモード](../integrations/delegate/) を参照。

## Delegate Edition — クイックスタート

### 前提条件

- **Git ≥ 2.41**
- サブスクリプション LLM 付き AI コーディングエージェント（Cursor、Claude Code、Codex 等）

### ステップ 1 — ビルドと Skill インストール

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

確認：

```bash
ocr-delegate -h
```

### ステップ 2 — Preview とレビュー

任意の Git リポジトリで：

```bash
cd path/to/your-repo
ocr-delegate preview
ocr-delegate rule <preview-出力のパス>
```

preview の mode/ref に基づき git で diff を取得。ホストエージェントがサブスクリプション LLM でレビュー — [デリゲーションモード](../integrations/delegate/) を参照。

---

## 上流 `ocr` CLI（本 Fork では未提供）

> 以下は **上流 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)** のみに適用されます。

### 前提条件

- **Git ≥ 2.41**
- **Node.js ≥ 18**
- **LLM API key**

### ステップ 1 —— CLI をインストールする

上流リポジトリからインストール — [上流インストール](https://github.com/alibaba/open-code-review#installation) を参照（npm グローバルパッケージ `@alibaba-group/open-code-review`）。

```bash
ocr version
```

> その他の方法は [インストール](../installation/) を参照してください。

### ステップ 2 —— LLM を設定する

```bash
ocr config provider
```

組み込みまたはカスタムの provider を選択し、API key を入力し、model を選び、設定ファイルに保存したうえで `ocr llm test` を 1 回実行してエンドポイントを検証します。

```bash
ocr config model
```

### 代替方法：非インタラクティブコマンド

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### ステップ 3 —— 接続性をテストする

```bash
ocr llm test
```

### ステップ 4 —— 初回のレビューを実行する

```bash
cd path/to/your-repo
ocr review
ocr review --from main --to feature-branch
ocr review --commit abc123
```

> 完全な `ocr review` 引数は [CLI リファレンス](../cli-reference/) を参照。

#### 先にプレビュー

```bash
ocr review --preview
ocr review -c abc123 --preview
```

#### JSON 出力

```bash
ocr review --format json --audience agent > review.json
```

## 関連項目

- [デリゲーションモード](../integrations/delegate/) — 本 Fork の主ワークフロー。
- [インストール](../installation/) — すべてのインストール方法と OCR の状態ディレクトリ。
- [設定](../configuration/) — 各環境変数、config key、組み込み provider。
- [CLI リファレンス](../cli-reference/) — 各サブコマンド、引数、出力モード。
- [レビュールール](../review-rules/) — レビュー内容をカスタマイズします。
- [インテグレーション](../integrations/agent-skill/) — OCR を Claude Code、Agent skill、CI に組み込みます。
- [FAQ](../faq/) — 既知のエラーと対策。
