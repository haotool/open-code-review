<div align="center">
  <h1>Open Code Review — Delegate Edition</h1>
  <p><strong>A security-hardened fork by <a href="https://github.com/haotool">haotool</a></strong></p>
</div>

<p align="center">
  <a href="README.zh-CN.md">简体中文</a> · <a href="README.zh-TW.md">繁體中文</a> · <a href="README.ja-JP.md">日本語</a> · <a href="README.ko-KR.md">한국어</a> · <a href="README.ru-RU.md">Русский</a>
</p>

<p align="center">
  <a href="https://github.com/haotool/open-code-review-delegate/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/haotool/open-code-review-delegate/ci.yml?style=flat-square" /></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/haotool/open-code-review-delegate?style=flat-square" /></a>
  <a href="https://github.com/alibaba/open-code-review"><img alt="Upstream" src="https://img.shields.io/badge/upstream-alibaba%2Fopen--code--review-blue?style=flat-square" /></a>
</p>

---

## Important Notice

This repository is **not** an official Alibaba product. It is an independent fork maintained by [haotool](https://github.com/haotool), derived from [alibaba/open-code-review](https://github.com/alibaba/open-code-review) under the Apache License 2.0.

- Repository: [`haotool/open-code-review-delegate`](https://github.com/haotool/open-code-review-delegate) (requests to the former `haotool/open-code-review` URL redirect automatically).
- We do **not** use Alibaba trademarks or imply official endorsement.
- See [NOTICE](NOTICE) for upstream copyright attribution and a summary of fork modifications.
- See [LICENSE](LICENSE) for the full Apache-2.0 terms.

## What Is This?

Open Code Review (Delegate Edition) ships a single CLI binary — **`ocr-delegate`** — that provides deterministic code-review engineering:

- **File selection** — precisely determines which changed files need review
- **Rule resolution** — matches review rules to each file from `rule.json`
- **Zero network** — the binary never calls an LLM or opens outbound connections
- **Zero credentials** — no API keys, tokens, or provider configuration required

Your AI coding agent (Cursor, Claude Code, Codex, etc.) performs the actual review using its own subscription LLM. `ocr-delegate` handles the engineering; the host agent handles the intelligence.

This design eliminates supply-chain risks from npm auto-updaters, prebuilt binary downloads, and embedded provider endpoints present in the upstream distribution model.

## Quick Install (3 Steps)

Build from source — no npm, no prebuilt binaries, no API keys:

```bash
# 1. Clone
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate

# 2. Build ocr-delegate
make build

# 3. Install the agent skill
make install-skill
# Claude Code users:
make install-skill SKILL_DIR=~/.claude/skills
```

Verify:

```bash
which ocr-delegate && ocr-delegate -h
```

Installation verified by `test/e2e/dryrun.sh` (see [E2E verification](test/e2e/README.md)).

## Usage

The skill at `skills/open-code-review/SKILL.md` is the single source of truth for the delegation workflow. Summary:

### 1. Preview — determine scope

```bash
ocr-delegate preview [--from main --to feature] [--commit <hash>] [-b "context"]
```

### 2. Get rules for files

```bash
ocr-delegate rule <path1> <path2> ...
```

### 3. Get diffs and review

Use git based on mode/ref metadata from preview output. The host agent reviews each file and produces structured findings.

See the [skill documentation](skills/open-code-review/SKILL.md) for the complete workflow, security discipline (T7/T9), and output schema.

See the [agent support matrix](docs/AGENT_SUPPORT.md) for the formally supported Codex, Claude Code, and Cursor installation paths.

## Security Posture

| Property | Guarantee |
|----------|-----------|
| Zero outbound | `ocr-delegate` makes no network connections (verified by CI string scan + E2E) |
| Zero credentials | No LLM provider config needed for delegate mode |
| Minimal dependency closure | Forbidden modules (`internal/llm`, `telemetry`, `mcp`, `viewer`, `session`) excluded from binary |
| Source-first distribution | No npm wrapper, no prebuilt binary downloads, no auto-updater |
| Supply-chain stripped | Upstream install scripts, npm packages, and GitHub Action removed |

Full details: [SECURITY.md](SECURITY.md)

## Build Reproduction & Verification

This fork is **source-first**: releases do not ship prebuilt binaries. Reproduce locally and verify integrity yourself:

```bash
make build                          # produces dist/ocr-delegate
shasum -a 256 dist/ocr-delegate     # local checksum for your build

# Cross-platform local builds (optional)
make dist                           # builds all platforms + sha256sum.txt
```

## Development

```bash
make build       # build ocr-delegate
make test        # run test suite
make check       # fmt + vet + mod tidy
make coverage    # coverage report (80% threshold)
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Upstream Legacy Directories

The following directories are retained from upstream but **not maintained** in this fork:

| Directory | Status |
|-----------|--------|
| `pages/` | Upstream landing page — not built or deployed by this fork |
| `extensions/` | Upstream VS Code extension — requires full `ocr` CLI, not `ocr-delegate` |

## License & Attribution

Licensed under [Apache-2.0](LICENSE). Copyright notices for upstream Alibaba contributors are preserved in [NOTICE](NOTICE).

This fork modifies upstream substantially. See [CHANGELOG.md](CHANGELOG.md) for v1.0.0 changes.

## Release Checklist (S1–S7)

| ID | Requirement | Verification |
|----|-------------|--------------|
| S1 | Zero outbound network | `test/e2e/dryrun.sh` + CI Gate 4 string scan |
| S2 | Zero credentials | E2E dry-run + delegate-only binary |
| S3 | Deterministic engine | 26-scenario golden equivalence suite |
| S4 | Source-first + delegate-only skill | No npm/postinstall supply chain in install paths; skill references `ocr-delegate` only (not full `ocr` CLI or external LLM API) |
| S5 | Install ≤3 steps | E2E dry-run three-step verification |
| S6 | Test suite passes | CI Gate 2 (`go test ./...`) |
| S7 | Injection resistance | Adversarial fixture library + E2E spot-checks |
