---
title: Agent Skill
sidebar:
  order: 1
---

> **Delegate Edition (this fork):** Install the skill with **`make install-skill`** after building **`ocr-delegate`** from source. See [Delegation Mode](../delegate/) for the canonical workflow.

## Delegate Edition — agent skill

This fork ships a delegate-first SKILL at
[`skills/open-code-review/SKILL.md`](https://github.com/haotool/open-code-review-delegate/blob/main/skills/open-code-review/SKILL.md).
It uses **`ocr-delegate`** for file selection and rule resolution; your host agent
performs the review with its subscription LLM.

### Prerequisites

- **Git ≥ 2.41**
- An AI coding agent with a subscription LLM (Cursor, Claude Code, Codex, etc.)

### Install

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

Verify:

```bash
ocr-delegate -h
which ocr-delegate
```

Claude Code users can target a custom skills directory:

```bash
make install-skill SKILL_DIR=~/.claude/skills
```

### Workflow

The skill drives a five-step loop — engineering scaffolding via `ocr-delegate`, reasoning via the host agent:

1. **Preview** — `ocr-delegate preview` lists reviewable files plus mode/ref metadata.
2. **Rules** — `ocr-delegate rule <paths…>` resolves review rules grouped by content.
3. **Diffs** — use git based on preview mode/ref metadata (see [Delegation Mode](../delegate/)).
4. **Review** — the host agent reviews each file with its subscription LLM.
5. **Report** — classify findings (Critical/High/Medium/Low) per the skill schema.

Example in any Git repository:

```bash
cd path/to/your-repo
ocr-delegate preview -b "PR context"
ocr-delegate rule internal/api/handler.go
```

Full flag reference and security discipline (T7/T9) live in the skill manifest and [Delegation Mode](../delegate/).

### Plugin wrappers

OpenCode, Claude Code, and Codex plugin entry points under `plugins/open-code-review/` point at the same skill. See [Claude Code Plugin](../claude-code/) for slash-command usage.

## See Also

- [Delegation Mode](../delegate/) — primary workflow for this fork.
- [QuickStart](../../quickstart/) — build, install-skill, first preview.
- [Installation](../../installation/) — prerequisites and skill install paths.

---

## Upstream agent skill (not shipped in this fork)

> The following applies to **upstream [alibaba/open-code-review](https://github.com/alibaba/open-code-review)** only.

Upstream ships a self-contained skill that invokes the full `ocr` CLI (including automatic npm-based CLI install on first run). See the upstream repository for `npx skills add`, Anthropic Agent SDK integration, and the embedded LLM review workflow.
