# Open Code Review — Plugin Pack (Delegate Edition)

> **Fork notice:** This repository is maintained by [haotool](https://github.com/haotool), not Alibaba.
> Plugins point at the **Delegate Edition** skill and `ocr-delegate` CLI — not upstream `ocr` + npm.

## Layout

| Path | Purpose |
|------|---------|
| `.claude-plugin/marketplace.json` | Local Codex/Claude marketplace catalog |
| `skills/open-code-review/` | Complete bundled delegate skill; checked against the repo-root SSOT |
| `skills/open-code-review-delegate/` | Complete bundled alias skill for legacy plugin IDs |
| `.codex-plugin/plugin.json` | Codex plugin manifest |
| `.cursor-plugin/plugin.json` | Cursor plugin manifest |
| `claude-code/commands/delegate-review.md` | Claude slash command for delegate workflow |
| `claude-code/skills/` | Complete Claude Code skill package |

## Install (recommended)

Build from source and install the skill — no npm:

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
make install-skill
# Claude Code:
make install-skill SKILL_DIR=~/.claude/skills
```

## Claude Code

Copy or symlink `claude-code/` into your Claude plugin path, or use the repo-root skill directly.

## Codex / Cursor

For Codex, add this checkout as a local marketplace and install the Delegate
plugin:

```bash
codex plugin marketplace add /path/to/open-code-review
codex plugin add open-code-review-delegate-codex@open-code-review
```

For Cursor, install the plugin directory using its local plugin flow. The
Codex and Cursor manifests point to the complete bundled skills; they do not
depend on a checkout-relative repo-root file.

See [the support matrix](../../docs/AGENT_SUPPORT.md) for the platform contract
and the optional live host-agent smoke test.

## Upstream

For the full `ocr` CLI + npm distribution, see [alibaba/open-code-review](https://github.com/alibaba/open-code-review).
