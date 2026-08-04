---
title: Agent Skill
sidebar:
  order: 1
---

> **Delegate Edition（本 Fork）：** 从源码构建 **`ocr-delegate`** 后，用 **`make install-skill`** 安装 Skill。标准流程见 [委托模式](../delegate/)。

## Delegate Edition — Agent Skill

本 Fork 在
[`skills/open-code-review/SKILL.md`](https://github.com/haotool/open-code-review-delegate/blob/main/skills/open-code-review/SKILL.md)
提供委托优先的 SKILL。它用 **`ocr-delegate`** 做文件筛选与规则解析；宿主 Agent 使用其订阅 LLM 完成审查。

### 前置条件

- **Git ≥ 2.41**
- 带订阅 LLM 的 AI 编码代理（Cursor、Claude Code、Codex 等）

### 安装

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

验证：

```bash
ocr-delegate -h
which ocr-delegate
```

Claude Code 用户可指定 Skill 目录：

```bash
make install-skill SKILL_DIR=~/.claude/skills
```

### 工作流

Skill 驱动五步循环 — 工程脚手架由 `ocr-delegate` 完成，推理由宿主 Agent 完成：

1. **Preview** — `ocr-delegate preview` 列出待审文件及 mode/ref 元数据。
2. **Rules** — `ocr-delegate rule <paths…>` 按内容分组解析审查规则。
3. **Diffs** — 根据 preview 的 mode/ref 使用 git 获取 diff（见 [委托模式](../delegate/)）。
4. **Review** — 宿主 Agent 用订阅 LLM 逐文件审查。
5. **Report** — 按 Skill 模式对发现项分级（Critical/High/Medium/Low）。

示例（任意 Git 仓库）：

```bash
cd path/to/your-repo
ocr-delegate preview -b "PR 上下文"
ocr-delegate rule internal/api/handler.go
```

完整参数与安全纪律（T7/T9）见 Skill 清单与 [委托模式](../delegate/)。

### 插件封装

`plugins/open-code-review/` 下的 OpenCode、Claude Code、Codex 插件入口指向同一 Skill。Slash 命令用法见 [Claude Code 插件](../claude-code/)。

## 参见

- [委托模式](../delegate/) — 本 Fork 的标准工作流。
- [快速开始](../../quickstart/) — 构建、install-skill、首次 preview。
- [安装](../../installation/) — 前置条件与 Skill 安装路径。

---

## 上游 Agent Skill（本 Fork 未提供）

> 以下内容仅适用于 **上游 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)**。

上游提供自包含 Skill，调用完整 `ocr` CLI（含首次运行时的 npm 自动安装）。`npx skills add`、Anthropic Agent SDK 集成及内嵌 LLM 审查流程见上游仓库文档。
