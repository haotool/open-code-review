<div align="center">
  <h1>Open Code Review — Delegate Edition</h1>
  <p><strong>由 <a href="https://github.com/haotool">haotool</a> 维护的安全加固 Fork</strong></p>
</div>

<p align="center">
  <a href="README.md">English</a> | 简体中文 | <a href="README.zh-TW.md">繁體中文</a> | <a href="README.ja-JP.md">日本語</a> | <a href="README.ko-KR.md">한국어</a> | <a href="README.ru-RU.md">Русский</a>
</p>

<p align="center">
  <a href="https://github.com/haotool/open-code-review-delegate/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/haotool/open-code-review-delegate/ci.yml?style=flat-square" /></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/haotool/open-code-review-delegate?style=flat-square" /></a>
  <a href="https://github.com/alibaba/open-code-review"><img alt="Upstream" src="https://img.shields.io/badge/upstream-alibaba%2Fopen--code--review-blue?style=flat-square" /></a>
</p>

---

## 重要说明

本仓库**不是**阿里巴巴官方产品。它是 [haotool](https://github.com/haotool) 在 Apache License 2.0 下独立维护的 Fork，源自 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)。

- 仓库地址：[`haotool/open-code-review-delegate`](https://github.com/haotool/open-code-review-delegate)（原 `haotool/open-code-review` URL 会自动重定向）。
- 不使用阿里巴巴商标，也不暗示官方背书。
- 上游版权说明见 [NOTICE](NOTICE)。
- 完整许可条款见 [LICENSE](LICENSE)。

## 这是什么？

Open Code Review（Delegate Edition）提供单一 CLI 二进制 **`ocr-delegate`**，负责确定性的代码审查工程：

- **文件筛选** — 精准判断哪些变更文件需要审查
- **规则解析** — 依 `rule.json` 为每个文件匹配审查规则
- **零外连** — 二进制不调用 LLM、不建立 outbound 连接
- **零凭据** — 无需 API Key、Token 或 provider 配置

您的 AI 编程助手（Cursor、Claude Code、Codex 等）使用**订阅制 LLM** 执行实际审查；`ocr-delegate` 负责工程面，宿主 Agent 负责智能面。

此设计移除了上游发行模式中 npm 自动更新、预编译二进制下载与内置 provider 端点等供应链风险。

## 快速安装（3 步）

从源码构建 — 无 npm、无预编译二进制、无 API Key：

```bash
# 1. Clone
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate

# 2. 构建 ocr-delegate
make build

# 3. 安装 Agent Skill
make install-skill
# Claude Code 用户：
make install-skill SKILL_DIR=~/.claude/skills
```

验证：

```bash
which ocr-delegate && ocr-delegate -h
```

安装流程由 `test/e2e/dryrun.sh` 验证（见 [E2E 验证](test/e2e/README.md)）。

## 用法

`skills/open-code-review/SKILL.md` 为委派工作流的单一真相来源（SSOT）。摘要：

### 1. Preview — 确定范围

```bash
ocr-delegate preview [--from main --to feature] [--commit <hash>] [-b "context"]
```

### 2. 获取文件规则

```bash
ocr-delegate rule <path1> <path2> ...
```

### 3. 获取 diff 并审查

依 preview 输出的 mode/ref 使用 git；宿主 Agent 审查各文件并产出结构化 findings。

完整流程、安全纪律（T7/T9）与输出 schema 见 [Skill 文档](skills/open-code-review/SKILL.md)。

正式支持的 Codex、Claude Code、Cursor 安装路径见 [Agent 支持矩阵](docs/AGENT_SUPPORT.md)。

## 安全 posture

| 属性 | 保证 |
|------|------|
| 零外连 | `ocr-delegate` 不建立网络连接（CI 字符串扫描 + E2E 验证） |
| 零凭据 | Delegate 模式无需 LLM provider 配置 |
| 最小依赖闭包 | 二进制排除 `internal/llm`、`telemetry`、`mcp`、`viewer`、`session` |
| 源码优先 | 无 npm 包装、无预编译下载、无自动更新 |
| 供应链精简 | 已移除上游 install 脚本、npm 包与 GitHub Action |

详情：[SECURITY.md](SECURITY.md)

## 构建复现与验证

本 Fork 为**源码优先**：不发布预编译二进制。请在本机重现并自行验证完整性：

```bash
make build                          # 产出 dist/ocr-delegate
shasum -a 256 dist/ocr-delegate     # 本机 checksum

# 跨平台本机构建（可选）
make dist                           # 全平台 + sha256sum.txt
```

## 开发

```bash
make build       # 构建 ocr-delegate
make test        # 测试套件
make check       # fmt + vet + mod tidy
make coverage    # 覆盖率报告（80% 门槛）
```

见 [CONTRIBUTING.zh-CN.md](CONTRIBUTING.zh-CN.md)。

## 上游遗留目录

以下目录保留自上游，但**本 Fork 不维护**：

| 目录 | 状态 |
|------|------|
| `pages/` | 上游 landing — 本 Fork 不构建/部署 |
| `extensions/` | 上游 VS Code 扩展 — 需完整 `ocr` CLI，非 `ocr-delegate` |

## 许可与归属

[Apache-2.0](LICENSE) 授权。上游 Alibaba 贡献者版权声明保留于 [NOTICE](NOTICE)。

本 Fork 对上游有大量修改。v1.0.0 变更见 [CHANGELOG.md](CHANGELOG.md)。

## 发布检查清单 (S1–S7)

| ID | 要求 | 验证 |
|----|------|------|
| S1 | 零外连 | E2E + CI Gate 4 字符串扫描 |
| S2 | 零凭据 | E2E dry-run + delegate-only 二进制 |
| S3 | 确定性引擎 | 26 场景 golden equivalence |
| S4 | 源码优先 + delegate-only skill | 安装路径无 npm/postinstall 供应链；Skill 仅引用 `ocr-delegate`（非完整 `ocr` CLI / 外部 LLM API） |
| S5 | 安装 ≤3 步 | E2E 三步验证 |
| S6 | 测试通过 | CI Gate 2（`go test ./...`） |
| S7 | 注入抵抗 | 对抗 fixture + E2E spot-check |
