---
title: 委托模式
sidebar:
  order: 5
---

> **Delegate Edition（haotool fork）：** 本仓库的标准工作流。从源码构建 `ocr-delegate` — 无 npm、OCR 端无需 LLM API Key。

`ocr-delegate` 负责确定性工程（文件筛选、规则解析），宿主 Agent（Cursor、Claude Code、Codex 等）使用自身订阅 LLM 执行实际代码审查。二进制不会发起出站网络连接。

## 何时使用委托模式

适用于已内置 LLM 订阅的 AI 编码代理。无需单独配置模型端点，直接复用宿主 Agent 的额度完成审查推理。

适用于：

1. 订阅制 AI 编码代理，希望复用已有额度进行代码审查。
2. 只需 OCR 的工程脚手架（文件过滤、规则解析、排除逻辑），LLM 推理由宿主 Agent 完成。
3. 构建自定义 Agent 流水线，需要结构化输入（文件列表 + 规则）。

## 前置条件

从源码构建（见 [README](https://github.com/haotool/open-code-review-delegate#quick-install-3-steps)）：

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
ocr-delegate -h
```

OCR 端无需 LLM 配置或环境变量。

## 安装 Skill

```bash
make install-skill
# Claude Code 用户：
make install-skill SKILL_DIR=~/.claude/skills
```

主 Skill 位于 `skills/open-code-review/SKILL.md`；兼容别名位于 `skills/open-code-review-delegate/SKILL.md`。`plugins/open-code-review/` 下的插件包装指向同一 Skill。

## 工作流程

### 第 1 步：Preview — 确定审查范围

```bash
ocr-delegate preview [--from <ref> --to <ref>] [--commit <hash>] [--exclude <patterns>]
```

### 第 2 步：获取文件规则

```bash
ocr-delegate rule <path1> <path2> ...
```

### 第 3 步：获取 diff

根据 Step 1 的 mode/ref 信息，直接使用 git：

**Range 模式**（提供了 merge\_base）：

```bash
git diff <merge_base>..<to> -- <path>
```

**Commit 模式**：

```bash
git show <commit> -- <path>
```

**Workspace 模式**：

```bash
git diff HEAD -- <path>        # 已跟踪文件
cat <path>                     # 新建的未跟踪文件
```

### 第 4 步：逐文件审查

结合 Step 2 的规则组，由宿主 Agent 的 LLM 完成审查。

### 第 5 步：输出报告

按 Skill  schema 标注严重级别。

## 子命令

| 命令 | 用途 |
|------|------|
| `ocr-delegate preview` | 列出可审查文件与 mode/ref 元数据 |
| `ocr-delegate rule <path...>` | 按规则内容分组解析 |

## 共享参数

| 参数 | 说明 |
|------|------|
| `--from <ref>` | Range 模式的源 ref |
| `--to <ref>` | Range 模式的目标 ref |
| `-c, --commit <hash>` | 单 commit 模式 |
| `--repo <path>` | 仓库根目录（默认：cwd） |
| `--rule <path>` | 自定义 rule.json 路径 |
| `--exclude <patterns>` | 逗号分隔的排除模式 |
| `-b, --background <text>` | 业务上下文 |
| `-B, --background-file <path>` | 从 Markdown 文件读取业务上下文 |

## 与上游对比

| 模式 | 谁调用 LLM？ | 本 fork？ |
|------|-------------|-----------|
| Agent Skill / `ocr review` | OCR 内置 LLM | 未提供 — 见带横幅的上游文档 |
| **委托（`ocr-delegate`）** | 宿主 Agent | **是 — 主工作流** |

## 参见

- [仓库 README](https://github.com/haotool/open-code-review-delegate)
- [E2E dry-run](https://github.com/haotool/open-code-review-delegate/tree/main/test/e2e)
