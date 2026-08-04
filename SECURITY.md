# Security Policy

## Fork Security Stance

This project is a security-hardened fork of [alibaba/open-code-review](https://github.com/alibaba/open-code-review). The delivered artifact is a delegate-only binary `ocr-delegate`: a deterministic engine with zero outbound network capability. The only code exfiltration point is the host agent's subscription LLM (a user-informed choice).

Threat identifiers (T1, T2, etc.) refer to the fork threat model documented below.

### Key Guarantees

1. **No npm wrapper or prebuilt download chain (T1)** — The upstream npm wrapper (`bin/ocr.js` with background `npm i -g` upgrades), postinstall prebuilt download chain (`scripts/install.js`, six-platform `npm/` packages), and `install.sh`/`install.ps1` download scripts have been removed. Only source builds (`go build -trimpath`) are supported; no prebuilt binaries are published or downloaded.

2. **FileReader absolute-path defect: upstream code unchanged (T6)** — The upstream FileReader does not hard-reject absolute paths. Per the fork's derivative decision, this upstream code is not modified; the delivered `ocr-delegate` dependency closure excludes that path (delegate does not register the `file_read` tool). Closure is verified by CI Gate 3 (`go list -deps ./cmd/ocr-delegate`).

3. **Telemetry, MCP, viewer, session excluded from delivery (T2/T4/T5)** — These modules and `internal/llm` (including LLM provider endpoint static tables) are not in the `ocr-delegate` dependency closure: no telemetry export, no listeners, no session JSONL persistence.

4. **Upstream GitHub Action unsupported (T8)** — Upstream `action.yml` and `ocr-review.yml` workflow have been removed. This fork's release workflow does not build or upload prebuilt binaries.

5. **Source-first distribution + delegate-only skill (S4)** — No npm wrapper, postinstall prebuilt download chain, or `@alibaba-group/open-code-review` install commands in canonical install paths. The agent skill references `ocr-delegate` only; it must not invoke the full `ocr review` CLI or configure external LLM API endpoints.

6. **Injection resistance at skill boundary (T7)** — Prompt injection defense is enforced by the host agent following skill security discipline. Adversarial fixtures in `testdata/adversarial/` validate expected behavior.

7. **Rule tampering detection (T9)** — Changes to `.opencodereview/rule.json` are flagged as high-priority review targets per skill workflow.

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

Only the latest released version receives security updates.

## Reporting a Vulnerability

**Please do NOT report security vulnerabilities through public GitHub issues.**

Use **GitHub Private Vulnerability Reporting**: [Security Advisories](https://github.com/haotool/open-code-review-delegate/security/advisories/new)

### What to Include

- Description of the vulnerability and potential impact
- Step-by-step reproduction instructions
- Affected version(s)
- Suggested fix or mitigation, if available

## Response Timeline

- **Acknowledgment**: within 3 business days
- **Initial Assessment**: within 7 business days
- **Fix & Disclosure**: within 14 days for confirmed critical/high-severity issues

## Scope

**In scope:**

- Remote code execution or command injection via crafted diffs, configs, or LLM responses
- Credential or API key leakage through logs, telemetry, or output files
- Path traversal allowing reads/writes outside the intended working directory
- Exploitable vulnerabilities in dependencies used by `ocr-delegate`

**Out of scope:**

- Issues in third-party LLM providers or APIs
- Denial-of-service requiring local access
- Social engineering attacks
- Host agent behavior when not following skill security discipline

## Releases

This fork uses source-first distribution. Obtain source from GitHub Releases and build locally:

```bash
make build
shasum -a 256 dist/ocr-delegate
```

## Recognition

We appreciate responsible disclosure. Reporters will be credited in release notes unless they prefer anonymity.
