# Changelog

All notable changes to this fork are documented here.

This project adheres to [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.0.0] - 2026-08-01

First release of the haotool Delegate Edition — a security-hardened fork of [alibaba/open-code-review](https://github.com/alibaba/open-code-review).

Epic and issue traceability: #1 (epic), #2–#11 (completed work items).

### Added

- **`ocr-delegate` CLI** — deterministic file selection and rule resolution without LLM calls (#3)
- **`internal/delegatecli`** — single source of truth for delegate preview/rule commands (#3, #4)
- **`internal/preview`** — extracted LLM-free file selection engine (#2)
- **Makefile targets** — `build`, `build-full`, `install-skill`, cross-platform `dist`, `sha256sum` (#5)
- **Four-gate CI** — build all packages, test suite, dependency closure check, provider string scan (#5)
- **Delegation skill** — `skills/open-code-review/SKILL.md` with T7/T9 security discipline (#7)
- **Golden equivalence suite** — 26 scenarios validating delegate CLI output parity (#8)
- **Adversarial fixture library** — T7/T9 samples for injection and rule-tampering regression (#9)
- **E2E dry-run** — `test/e2e/dryrun.sh` for install journey and zero-outbound verification (#10)
- **NOTICE** — Apache-2.0 upstream attribution and fork modification summary (#11)
- **Release checklist** — S1–S7 verification mapping in README (#11)
- **i18n README/CONTRIBUTING** — restored with fork disclaimers (ja/ko/ru/zh)

### Changed

- **`cmd/opencodereview`** — delegate subcommands converge on `internal/delegatecli` (#4)
- **`internal/relocation`** — LLM re-location logic extracted from agent package (#2)
- **Release workflow** — source-first only; no prebuilt binary upload or npm publish (#6)
- **SECURITY.md** — fork-specific threat model (T1–T9) and reporting process (#6)
- **README** — haotool product story, 3-step install, security posture, fork attribution (#11)
- **CI workflows** — pages-ci, vscode-ext, translation-sync restored for fork maintenance
- **pages/** — retained for upstream sync via `pages-ci`; GitHub Pages deploy (`deploy-pages.yml`) intentionally omitted (see README)
- **plugins/** — thin wrappers pointing to delegate skill and `ocr-delegate`

### Removed

- npm distribution chain (`bin/ocr.js`, `npm/` packages, `scripts/install.js`, auto-updater) (#6)
- Install download scripts (`install.sh`, `install.ps1`) (#6)
- Upstream GitHub Action (`action.yml`) and review workflow (#6)
- Prebuilt binary upload from release pipeline (#6)

### Security

- Binary dependency closure verified: no `internal/llm`, `telemetry`, `mcp`, `viewer`, or `session` (#5)
- Provider endpoint domain strings absent from `ocr-delegate` binary (CI Gate 4) (#5)
- No credential or API key requirements for delegate workflow (#7)
- Supply-chain attack surface reduced by removing npm wrapper and prebuilt downloads (#6)

[1.0.0]: https://github.com/haotool/open-code-review-delegate/releases/tag/v1.0.0
