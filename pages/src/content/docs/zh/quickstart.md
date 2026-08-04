---
title: 快速开始
sidebar:
  order: 3
---

几分钟内跑通第一次代码评审。

> **Delegate Edition（本 Fork）：** 标准工作流使用从源码构建的 **`ocr-delegate`** — 无 npm、OCR 端无需 LLM API Key。完整流程见 [委托模式](../integrations/delegate/)。

## Delegate Edition — 快速开始

### 前置条件

- **Git ≥ 2.41**
- 带订阅 LLM 的 AI 编码代理（Cursor、Claude Code、Codex 等）

### 第 1 步 — 构建并安装 Skill

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

验证：

```bash
ocr-delegate -h
```

### 第 2 步 — Preview 并审查

在任意 Git 仓库中：

```bash
cd path/to/your-repo
ocr-delegate preview
ocr-delegate rule <preview-输出的路径>
```

根据 preview 的 mode/ref 使用 git 获取 diff；宿主 Agent 使用其订阅 LLM 完成审查 — 见 [委托模式](../integrations/delegate/)。

---

## 上游 `ocr` CLI（本 Fork 未提供）

> 以下内容仅适用于 **上游 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)**。

### 前置条件

- **Git ≥ 2.41**
- **Node.js ≥ 18**
- **LLM API key**

### 第 1 步 —— 安装 CLI

请在上游仓库安装 — 见 [上游安装说明](https://github.com/alibaba/open-code-review#installation)（npm 全局包 `@alibaba-group/open-code-review`）。

```bash
ocr version
```

> 更多方式见 [安装](../installation/)。

### 第 2 步 —— 配置 LLM

```bash
ocr config provider
```

它会让你选择一个内置或自定义 provider、填入 API key、挑选 model，保存到配置文件后自动运行一次 `ocr llm test` 验证端点。之后想换模型：

```bash
ocr config model
```

### 替代方式：非交互命令

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### 第 3 步 —— 测试连通性

```bash
ocr llm test
```

若出现 `no valid LLM endpoint configured`，请复查第 2 步配置。401 / 403 表示 token 错误或过期。

### 第 4 步 —— 运行首次评审

```bash
cd path/to/your-repo
ocr review
ocr review --from main --to feature-branch
ocr review --commit abc123
```

> 完整 `ocr review` 参数见 [CLI 参考](../cli-reference/)。

#### 先预览

```bash
ocr review --preview
ocr review -c abc123 --preview
```

#### JSON 输出

```bash
ocr review --format json --audience agent > review.json
```

## 相关

- [委托模式](../integrations/delegate/) — 本 Fork 主工作流。
- [安装](../installation/) — 所有安装方式与 OCR 状态目录。
- [配置](../configuration/) — 环境变量、配置项与内置 provider。
- [CLI 参考](../cli-reference/) — 子命令、参数与输出模式。
- [评审规则](../review-rules/) — 自定义审查范围。
- [集成](../integrations/agent-skill/) — 在 Claude Code、Agent Skill 或 CI 中使用。
- [FAQ](../faq/) — 常见错误与处理。
