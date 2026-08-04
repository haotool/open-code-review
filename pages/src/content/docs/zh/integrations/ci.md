---
title: CI/CD
sidebar:
  order: 4
---

> **Delegate Edition（本 Fork）：** 通过 npm 安装上游 `ocr` CLI 的 CI 流水线**不是**本仓库的标准工作流。请从源码构建 **`ocr-delegate`**，在本地通过 Agent Skill 完成审查 — 见 [委托模式](../delegate/)。

## Delegate Edition — CI/CD

本 Fork **不**提供基于 npm 的 CI 配方。`ocr-delegate` 是零外联的二进制，负责确定性文件筛选与规则解析；LLM 推理由宿主 Agent 使用其订阅配额完成。

### 推荐做法

在开发流程或 Agent 会话中运行代码审查：

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

在待审仓库中：

```bash
ocr-delegate preview --from origin/main --to HEAD
ocr-delegate rule <preview-输出的路径>
```

根据 preview 的 mode/ref 使用 git 获取 diff；宿主 Agent（Cursor、Claude Code、Codex 等）完成审查 — 见 [Skill 工作流](../delegate/)。

### 可选：自托管 Runner + 源码构建

若需自动化且不使用 npm，可在 Runner 上构建 `ocr-delegate` 并在 Job 脚本中调用 `preview` / `rule`。**不要**期望 `ocr-delegate` 单独产出 JSON 审查评论 — 后处理与 LLM 审查仍由宿主 Agent 负责。

GitHub Actions 示例：

```yaml
- name: Build ocr-delegate
  run: make build
- name: Preview changed files
  run: ./dist/ocr-delegate preview --from origin/${{ github.base_ref }} --to HEAD
```

按需扩展 Agent 流水线。

## 参见

- [委托模式](../delegate/) — preview → rule → 审查工作流。
- [快速开始](../../quickstart/) — 三步安装。
- [E2E dry-run](https://github.com/haotool/open-code-review-delegate/tree/main/test/e2e) — 零外联验证。

---

## 上游 CI/CD（本 Fork 未提供）

> 以下内容仅适用于 **上游 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)**。

上游提供 GitHub Actions 与 GitLab CI 流水线：安装 `ocr` CLI、从 Secrets 配置 LLM、运行 `ocr review --format json` 并回帖 PR/MR 行内评论。

- GitHub Actions：[`examples/github_actions/ocr-review.yml`](https://github.com/alibaba/open-code-review/blob/main/examples/github_actions/ocr-review.yml)
- GitLab CI：[`examples/gitlab_ci/.gitlab-ci.yml`](https://github.com/alibaba/open-code-review/blob/main/examples/gitlab_ci/.gitlab-ci.yml)

Secrets、定制与排障见上游文档。
