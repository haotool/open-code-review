# End-to-End Dry-Run Verification

Validates the haotool fork install journey and zero-outbound delegate workflow before release.

## What it checks

| ID | Requirement | Method |
|----|-------------|--------|
| S1 | Zero outbound network | `lsof -p` monitoring during `ocr-delegate` execution |
| S2 | Zero credentials | No LLM/API env vars required or read |
| S5 | Install ≤3 steps | `make build` → `make install-skill` after clone |
| T5 | No residue | Isolated `$HOME`; no new files under `~/.opencodereview/` |
| S4 | Source-first + delegate-only skill | No npm/postinstall supply chain in canonical install paths; skill references `ocr-delegate` only (not full `ocr` CLI or external LLM API) |
| Gate 5 | Removed distribution surface | `action.yml`, `bin/ocr.js`, and `ocr.js` cannot reappear |
| S7 | Adversarial spot-check | Samples 01 (injection) and 04 (rule weakening) |
| H1 | Live host-agent smoke | `test/e2e/host-agent-smoke.sh` with a caller-supplied agent command |

## Run

```bash
bash test/e2e/dryrun.sh
```

The deterministic dry-run does not claim that a subscription agent loaded the
skill. For a real host-agent check, provide the installed agent's local command
explicitly:

```bash
HOST_AGENT_CMD='<agent command that reads stdin and writes the review>' \
  bash test/e2e/host-agent-smoke.sh
```

The live smoke exits `77` (`NOT_RUN`) when no command is supplied. It never
turns an unavailable subscription agent into a passing test.

On success, `test/e2e/last-run-report.txt` is written with a summary (no absolute user paths).

## Scope

This script verifies the **OCR side** (binary behavior, install path, network isolation). The adversarial fixtures validate that injection-shaped data reaches the host-agent boundary safely; they do not prove a host agent obeyed the skill. Use the H1 live smoke for that claim — the binary itself remains deterministic.
