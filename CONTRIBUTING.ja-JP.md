# コントリビューション — Open Code Review (Delegate Edition)

[haotool](https://github.com/haotool) が [alibaba/open-code-review](https://github.com/alibaba/open-code-review)（Apache-2.0）から派生してメンテナンスする Fork への貢献に感謝します。

## セットアップ

```bash
git clone https://github.com/haotool/open-code-review-delegate.git
cd open-code-review-delegate
make build && make test
```

## 開発コマンド

```bash
make build && make build-full && make test && make check && make coverage
make install-skill SKILL_DIR=~/.cursor/skills
```

## ガイドライン

- 既存の Go 規約に従う
- `make check` を PR 前に実行
- 新機能にはテストを追加
- `ocr-delegate` の依存閉包を維持 — delegate パスから `internal/llm`、`telemetry`、`mcp`、`viewer`、`session` を import しない
- Conventional Commits（`feat(scope):`、`fix:`、`docs:`、`test:`、`ci:`、`chore:`）

## スコープ

- **`ocr-delegate`** が主要成果物。ネットワーク/資格情報面を拡大しない
- **`pages/`**、**`extensions/`** は上流レガシー — 本 Fork では非メンテ
- npm 配布、プリビルト DL、上流 LLM Action の再導入禁止

## ライセンス

貢献は [Apache-2.0](LICENSE) で提供されることに同意したものとみなします。
