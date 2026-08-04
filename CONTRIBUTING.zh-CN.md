# 参与贡献 — Open Code Review (Delegate Edition)

感谢关注！本 Fork 由 [haotool](https://github.com/haotool) 维护，源自 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)（Apache-2.0）。

## 行为准则

参与即表示您同意维护尊重、包容的社区环境。

## 贡献方式

- **报告 Bug** — 附带复现步骤
- **改进建议** — 提交 Feature Request
- **文档** — 修正笔误、补充示例
- **代码** — 修复 Bug、补充测试、改进 delegate CLI

## 入门

### 前置条件

- Go 1.25+
- Git >= 2.41
- Make

### 设置

```bash
git clone https://github.com/haotool/open-code-review-delegate.git
cd open-code-review-delegate
make build && make test
```

### 常用命令

```bash
make build       # 构建 ocr-delegate
make build-full  # 构建完整 opencodereview（仅开发）
make test        # 测试（含 race detector）
make check       # fmt + vet + mod tidy
make coverage    # 覆盖率（80% 阈值）
make install-skill SKILL_DIR=~/.cursor/skills
```

## 编码规范

- 遵循现有 Go 风格
- 提交前运行 `make check`
- 新行为需有测试
- 保持 `ocr-delegate` 依赖闭包干净 — delegate 路径不得导入 `internal/llm`、`telemetry`、`mcp`、`viewer`、`session`
- Commit 格式：`feat(scope):`、`fix(scope):`、`docs:`、`test:`、`ci:`、`chore:`

## Pull Request

1. Fork 仓库
2. 从 `main` 创建功能分支
3. 附带测试的改动
4. 运行 `make check && make test`
5. 提交 PR 并说明变更

## 范围说明

- **`ocr-delegate`** 是主要交付物；不得扩大其网络或凭据面
- **`pages/`**、**`extensions/`** 为上游遗留目录 — 本 Fork 不主动维护
- 不得重新引入 npm 分发、预编译下载或上游 LLM GitHub Action

## 许可

贡献即表示您同意在 [Apache-2.0](LICENSE) 下授权。
