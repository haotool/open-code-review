# Agent Support Matrix

Open Code Review Delegate Edition has one product contract: `ocr-delegate`
performs deterministic scope selection and rule resolution, while the host
coding agent performs the review with its own model. The delegate binary does
not call an LLM, read provider credentials, or open outbound connections.

## Supported platforms

| Platform | Skill support | Native plugin package | Installation path | Status |
|----------|---------------|-----------------------|-------------------|--------|
| Codex | Yes | `.claude-plugin/marketplace.json` + `.codex-plugin/plugin.json` | `codex plugin marketplace add <repo>` then `codex plugin add open-code-review-delegate-codex@open-code-review` | Package + H1 verified locally |
| Claude Code | Yes | `claude-code/.claude-plugin/plugin.json` + bundled `skills/` + `delegate-review` command | Install the `claude-code/` package and/or `make install-skill SKILL_DIR=~/.claude/skills` | Package-supported; H1 verified locally |
| Cursor | Yes | `.cursor-plugin/plugin.json` | Install the plugin directory or `make install-skill` | Package-supported; H1 pending |
| OpenCode | No | None | Upstream-only legacy files are not part of this product | Not supported |

All three supported platforms use the same bundled
`skills/open-code-review/SKILL.md`. The plugin package contains a complete copy
of both skill files so an installer does not need to resolve a path outside the
package.

## What “subscription LLM” means here

The host agent's normal subscription model performs the reasoning and emits
the review findings. No separate `OCR_LLM_*` endpoint or API key is required.
This is a skill/plugin integration contract, not an API adapter that can force
an arbitrary agent to load a skill.

The live host-agent smoke test is opt-in because each platform has a different
CLI. Run it with a trusted local command for the installed agent:

```bash
HOST_AGENT_CMD='<agent command that reads stdin and writes the review>' \
  bash test/e2e/host-agent-smoke.sh
```

Without `HOST_AGENT_CMD`, the test reports `NOT_RUN` and exits with status 77;
it must never be reported as a passing subscription-agent verification.

Current evidence: Codex loaded the installed plugin, completed a bounded
host-agent review, and returned `CODEX_SKILL_REVIEW_PASS`. Claude Code also
passed the live smoke with the nested `claude-code/` plugin package. Cursor
remains H1-pending until its native local CLI flow is run through the same
harness.

## Deliberately unsupported paths

- The upstream embedded-LLM `ocr review` workflow is not the Delegate Edition.
- npm installation, prebuilt download chains, and provider configuration are
  not supported installation paths.
- `pages/`, the upstream VS Code extension, and other legacy full-CLI surfaces
  are retained for upstream attribution only and are not delivery targets.

## Verification

- `bash test/ci/check-plugin-skills.sh` verifies complete skill copies and an
  isolated plugin package.
- `bash test/e2e/dryrun.sh` verifies the delegate binary and skill contract.
- `bash test/e2e/host-agent-smoke.sh` is the opt-in live host-agent check.
