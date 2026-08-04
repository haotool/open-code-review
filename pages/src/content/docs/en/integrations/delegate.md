---
title: Delegation Mode
sidebar:
  order: 5
---

> **Delegate Edition (haotool fork):** This is the canonical workflow for this repository.
> Build `ocr-delegate` from source — no npm, no LLM API keys on the OCR side.

`ocr-delegate` handles deterministic engineering (file selection, rule resolution)
while your host agent (Cursor, Claude Code, Codex, etc.) performs the actual code
review using its own subscription LLM. The binary never opens outbound network
connections.

## When to use delegation mode

Delegation mode is designed for subscription-based AI coding agents where you
already have an LLM bundled with the host agent. Instead of configuring a
separate model endpoint, you reuse the host agent's quota for review reasoning.

Use delegation mode when:

1. Your AI coding agent runs on a subscription plan and you want to reuse that
   quota for code review — no extra API key or model configuration needed.
2. You want OCR only for engineering scaffolding — file filtering, rule
   resolution, exclusion logic — while the host agent handles all LLM reasoning.
3. You are building a custom agent pipeline that needs structured inputs (file
   list + rules) for its own review step.

## Prerequisites

Build from source (see [README](https://github.com/haotool/open-code-review-delegate#quick-install-3-steps)):

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
ocr-delegate -h
```

No LLM configuration or environment variables are required on the OCR side.

## Install the skill

```bash
make install-skill
# Claude Code users:
make install-skill SKILL_DIR=~/.claude/skills
```

The canonical skill lives at `skills/open-code-review/SKILL.md`. A backward-compatible
alias is at `skills/open-code-review-delegate/SKILL.md`.

Plugin wrappers (OpenCode, Claude Code, Codex) point to the same skill under
`plugins/open-code-review/`.

## Workflow

### Step 1: Preview — determine what to review

```bash
ocr-delegate preview [--from <ref> --to <ref>] [--commit <hash>] [--exclude <patterns>]
```

Outputs:

- **mode** — workspace / range / commit
- **ref metadata** — from, to, commit, merge\_base
- **Reviewable file list** — paths, status, insertions/deletions
- **Excluded files** — with exclusion reason

Common invocations:

| Scenario | Command |
|----------|---------|
| Workspace changes | `ocr-delegate preview` |
| Branch comparison | `ocr-delegate preview --from main --to feature` |
| Single commit | `ocr-delegate preview -c abc123` |

### Step 2: Get rules for files

```bash
ocr-delegate rule <path1> <path2> ...
```

Pass the reviewable paths from Step 1. Output is grouped by rule content.

### Step 3: Get diffs

Use git directly, based on the mode/ref info from Step 1:

**Range mode** (merge\_base provided):

```bash
git diff <merge_base>..<to> -- <path>
```

**Commit mode**:

```bash
git show <commit> -- <path>
```

**Workspace mode**:

```bash
git diff HEAD -- <path>        # tracked files
cat <path>                     # new untracked files
```

### Step 4: Review each file

For each reviewable file:

1. Get its diff (Step 3)
2. Consult the matching Rule Group (Step 2) as the review checklist
3. Conduct a thorough review using the host agent's LLM

### Step 5: Report

Classify each finding by severity (Critical/High/Medium/Low) per the skill schema.

## Sub-commands reference

| Command | Purpose |
|---------|---------|
| `ocr-delegate preview` | List reviewable files + mode/ref metadata |
| `ocr-delegate rule <path...>` | Resolve review rules grouped by content |

## Shared flags

| Flag | Description |
|------|-------------|
| `--from <ref>` | Source ref for range mode |
| `--to <ref>` | Target ref for range mode |
| `-c, --commit <hash>` | Single commit mode |
| `--repo <path>` | Repository root (default: cwd) |
| `--rule <path>` | Custom rule.json path |
| `--exclude <patterns>` | Comma-separated exclude patterns |
| `-b, --background <text>` | Business context |
| `-B, --background-file <path>` | Business context from Markdown file |

## Upstream comparison

| Mode | Who calls the LLM? | In this fork? |
|------|-------------------|---------------|
| Agent Skill / `ocr review` | OCR (embedded LLM) | Not shipped — see upstream docs with banner |
| **Delegation (`ocr-delegate`)** | Host agent | **Yes — primary workflow** |

## See Also

- [Repository README](https://github.com/haotool/open-code-review-delegate) — install, security posture, E2E verification
- [E2E dry-run](https://github.com/haotool/open-code-review-delegate/tree/main/test/e2e) — zero-outbound checks
