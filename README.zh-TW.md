<div align="center">
  <h1>Open Code Review — Delegate Edition</h1>
  <p><strong>由 <a href="https://github.com/haotool">haotool</a> 維護的安全加固 Fork</strong></p>
</div>

<p align="center">
  <a href="README.md">English</a> · <a href="README.zh-CN.md">简体中文</a> · 繁體中文 · <a href="README.ja-JP.md">日本語</a> · <a href="README.ko-KR.md">한국어</a> · <a href="README.ru-RU.md">Русский</a>
</p>

<p align="center">
  <a href="https://github.com/haotool/open-code-review-delegate/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/haotool/open-code-review-delegate/ci.yml?style=flat-square" /></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/haotool/open-code-review-delegate?style=flat-square" /></a>
  <a href="https://github.com/alibaba/open-code-review"><img alt="Upstream" src="https://img.shields.io/badge/upstream-alibaba%2Fopen--code--review-blue?style=flat-square" /></a>
</p>

---

## 重要說明

本倉庫**不是**阿里巴巴官方產品。它是 [haotool](https://github.com/haotool) 在 Apache License 2.0 下獨立維護的 Fork，源自 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)。

- 倉庫位址：[`haotool/open-code-review-delegate`](https://github.com/haotool/open-code-review-delegate)（原 `haotool/open-code-review` URL 會自動重新導向）。
- 不使用阿里巴巴商標，也不暗示官方背書。
- 上游版權說明見 [NOTICE](NOTICE)。
- 完整授權條款見 [LICENSE](LICENSE)。

## 這是什麼？

Open Code Review（Delegate Edition）提供單一 CLI 二進位 **`ocr-delegate`**，負責確定性的程式碼審查工程：

- **檔案篩選** — 精準判斷哪些變更檔案需要審查
- **規則解析** — 依 `rule.json` 為每個檔案匹配審查規則
- **零外連** — 二進位不呼叫 LLM、不建立 outbound 連線
- **零憑證** — 無需 API Key、Token 或 provider 設定

您的 AI 程式助手（Cursor、Claude Code、Codex 等）使用**訂閱制 LLM** 執行實際審查；`ocr-delegate` 負責工程面，宿主 Agent 負責智慧面。

此設計移除了上游發行模式中 npm 自動更新、預編譯二進位下載與內建 provider 端點等供應鏈風險。

## 快速安裝（3 步）

從源碼建置 — 無 npm、無預編譯二進位、無 API Key：

```bash
# 1. Clone
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate

# 2. 建置 ocr-delegate
make build

# 3. 安裝 Agent Skill
make install-skill
# Claude Code 使用者：
make install-skill SKILL_DIR=~/.claude/skills
```

驗證：

```bash
which ocr-delegate && ocr-delegate -h
```

安裝流程由 `test/e2e/dryrun.sh` 驗證（見 [E2E 驗證](test/e2e/README.md)）。

## 用法

`skills/open-code-review/SKILL.md` 為委派工作流的單一真相來源（SSOT）。摘要：

### 1. Preview — 確定範圍

```bash
ocr-delegate preview [--from main --to feature] [--commit <hash>] [-b "context"]
```

### 2. 取得檔案規則

```bash
ocr-delegate rule <path1> <path2> ...
```

### 3. 取得 diff 並審查

依 preview 輸出的 mode/ref 使用 git；宿主 Agent 審查各檔並產出結構化 findings。

完整流程、安全紀律（T7/T9）與輸出 schema 見 [Skill 文件](skills/open-code-review/SKILL.md)。

正式支援的 Codex、Claude Code、Cursor 安裝路徑見 [Agent 支援矩陣](docs/AGENT_SUPPORT.md)。

## 安全姿態

| 屬性 | 保證 |
|------|------|
| 零外連 | `ocr-delegate` 不建立網路連線（CI 字串掃描 + E2E 驗證） |
| 零憑證 | Delegate 模式無需 LLM provider 設定 |
| 最小依賴閉包 | 二進位排除 `internal/llm`、`telemetry`、`mcp`、`viewer`、`session` |
| 源碼優先 | 無 npm 包裝、無預編譯下載、無自動更新 |
| 供應鏈精簡 | 已移除上游 install 腳本、npm 套件與 GitHub Action |

詳情：[SECURITY.md](SECURITY.md)

## 建置重現與驗證

本 Fork 為**源碼優先**：不發布預編譯二進位。請在本機重現並自行驗證完整性：

```bash
make build                          # 產出 dist/ocr-delegate
shasum -a 256 dist/ocr-delegate     # 本機 checksum

# 跨平台本機建置（可選）
make dist                           # 全平台 + sha256sum.txt
```

## 開發

```bash
make build       # 建置 ocr-delegate
make test        # 測試套件
make check       # fmt + vet + mod tidy
make coverage    # 覆蓋率報告（80% 門檻）
```

指南見 [CONTRIBUTING.md](CONTRIBUTING.md)（繁體說明亦可參考 [CONTRIBUTING.zh-CN.md](CONTRIBUTING.zh-CN.md)）。

## 上游遺留目錄

以下目錄保留自上游，但**本 Fork 不維護**：

| 目錄 | 狀態 |
|------|------|
| `pages/` | 上游 landing — 本 Fork 不建置/部署 |
| `extensions/` | 上游 VS Code 擴充 — 需完整 `ocr` CLI，非 `ocr-delegate` |

## 授權與歸屬

[Apache-2.0](LICENSE) 授權。上游 Alibaba 貢獻者版權聲明保留於 [NOTICE](NOTICE)。

本 Fork 對上游有大量修改。v1.0.0 變更見 [CHANGELOG.md](CHANGELOG.md)。

## 發布檢查清單 (S1–S7)

| ID | 要求 | 驗證 |
|----|------|------|
| S1 | 零外連 | E2E + CI Gate 4 字串掃描 |
| S2 | 零憑證 | E2E dry-run + delegate-only 二進位 |
| S3 | 確定性引擎 | 26 場景 golden equivalence |
| S4 | 源碼優先 + delegate-only skill | 安裝路徑無 npm/postinstall 供應鏈；Skill 僅引用 `ocr-delegate`（非完整 `ocr` CLI / 外部 LLM API） |
| S5 | 安裝 ≤3 步 | E2E 三步驗證 |
| S6 | 測試通過 | CI Gate 2（`go test ./...`） |
| S7 | 注入抵抗 | 對抗 fixture + E2E spot-check |
