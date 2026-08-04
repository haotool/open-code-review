---
title: QuickStart
sidebar:
  order: 3
---

Get your first code review running in a few minutes.

> **Delegate Edition (this fork):** The primary workflow uses **`ocr-delegate`** built from source — no npm, no LLM API keys on the OCR side. See [Delegation Mode](../integrations/delegate/) for the full workflow.

## Delegate Edition — quick start

### Prerequisites

- **Git ≥ 2.41**
- An AI coding agent with a subscription LLM (Cursor, Claude Code, Codex, etc.)

### Step 1 — Build and install the skill

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

Verify:

```bash
ocr-delegate -h
```

### Step 2 — Preview and review

In any Git repository:

```bash
cd path/to/your-repo
ocr-delegate preview
ocr-delegate rule <paths-from-preview>
```

Use git for diffs based on preview mode/ref metadata. Your host agent performs the review using its subscription LLM — see the [skill workflow](../integrations/delegate/).

---

## Upstream `ocr` CLI (not shipped in this fork)

> The following applies to **upstream [alibaba/open-code-review](https://github.com/alibaba/open-code-review)** only.

### Prerequisites

- **Git ≥ 2.41**
- **Node.js ≥ 18**
- **LLM API key**

### Step 1 — Install the CLI

Install from the upstream repository — see [upstream Installation](https://github.com/alibaba/open-code-review#installation) (npm global package `@alibaba-group/open-code-review`).

```bash
ocr version
```

> See [Installation](../installation/) for more methods.

### Step 2 — Configure an LLM

```bash
ocr config provider
```

It lets you pick a built-in or custom provider, enter an API key, choose a model, saves everything to the config file, and then runs `ocr llm test` once to verify the endpoint. To switch models later:

```bash
ocr config model
```

### Alternative: non-interactive command

In CI or a no-TUI environment, write to the same config directly with `ocr config set`:

```bash
ocr config set provider                    anthropic
ocr config set model                       claude-opus-4-6
ocr config set providers.anthropic.api_key sk-ant-xxxxxxxxxx
```

### Step 3 — Test connectivity

```bash
ocr llm test
```

If you get an error like `no valid LLM endpoint configured`, recheck the Step 2 config. A 401 / 403 means the token is wrong or expired.

### Step 4 — Run your first review

```bash
cd path/to/your-repo

# Workspace mode — reviews staged + unstaged + untracked changes (default)
ocr review

# Branch range — reviews `main..feature-branch`
ocr review --from main --to feature-branch

# Single commit — reviews the diff that commit introduced
ocr review --commit abc123
```

> See [CLI Reference](../cli-reference/) for the complete list of `ocr review` flags.

#### Preview first

```bash
ocr review --preview              # workspace
ocr review -c abc123 --preview    # commit
```

#### JSON output for systems

```bash
ocr review --format json --audience agent > review.json
```

## See Also

- [Delegation Mode](../integrations/delegate/) — primary workflow for this fork.
- [Installation](../installation/) — every install method and OCR's state directory.
- [Configuration](../configuration/) — every env var, config key, and built-in provider.
- [CLI Reference](../cli-reference/) — every sub-command, flag, and output mode.
- [Review Rules](../review-rules/) — customize what gets reviewed.
- [Integrations](../integrations/agent-skill/) — embed OCR in Claude Code, an Agent skill, or CI.
- [FAQ](../faq/) — known errors and remedies.
