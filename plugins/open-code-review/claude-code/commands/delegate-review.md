# Delegate code review (ocr-delegate)

Run a **delegate-mode** review: `ocr-delegate` for scope/rules, host agent for LLM review.

## Steps

1. `ocr-delegate preview` — determine files and mode
2. `ocr-delegate rule <paths...>` — resolve rules per file
3. Use git diff per preview metadata; produce structured findings (see skill)

## Install

```bash
make build && make install-skill SKILL_DIR=~/.claude/skills
```

Full workflow: `../skills/open-code-review/SKILL.md` in this plugin package.

**Not applicable:** upstream `ocr review` with embedded LLM — use delegate workflow instead.
