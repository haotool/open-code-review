# Security Assurance Case

> **Scope: Delegate Edition (haotool fork)**
> This assurance case applies to the delivered artifact **`ocr-delegate`** — a zero-network CLI plus the agent skill workflow. Threat identifiers **T1–T9** match [SECURITY.md](./SECURITY.md).
> Sections marked **(upstream only)** describe the upstream `opencodereview` / `ocr` product (embedded LLM, viewer, npm distribution) and are **not** claims for this fork.

## Threat Model (Delegate Edition)

### System Description

In Delegate Edition, `ocr-delegate`:

1. Reads git diff output from a local repository.
2. Resolves review rules from local JSON configuration.
3. Emits preview/rule output to stdout for the **host agent** to consume.
4. Makes **no outbound network calls** and holds **no LLM credentials**.

The host agent (Cursor, Claude Code, etc.) performs the actual review using its subscription LLM — a user-informed choice outside the binary's trust boundary.

### Actors

| Actor | Trust Level |
|-------|-------------|
| Local user | Trusted — invokes the CLI with full control over configuration |
| Host agent + subscription LLM | Semi-trusted — must follow skill security discipline (T7) |
| Git repository | Semi-trusted — diffs may contain adversarial content |
| Network | **Not used by `ocr-delegate`** |

### Trust Boundaries

```
┌──────────────────────────────────────────────────┐
│  User Machine (Trusted Zone)                     │
│                                                  │
│  ┌──────────┐    ┌──────────────┐    ┌────────┐  │
│  │ Git Repo │───▶│ ocr-delegate │───▶│ stdout │  │
│  │ (diffs)  │    │ (deterministic)│   │ (data) │  │
│  └──────────┘    └──────────────┘    └────────┘  │
│        Trust ────────▶│                           │
│        Boundary 1     │                           │
│                       ▼                           │
│              ┌─────────────────┐                  │
│              │  Host Agent     │── Boundary 2     │
│              │  (skill + LLM)  │   (user choice)  │
│              └─────────────────┘                  │
└──────────────────────────────────────────────────┘
```

1. **Git → CLI**: Diff content may contain crafted payloads. Parsed with strict format validation; only `git` is invoked with hardcoded subcommands.
2. **CLI → Host agent**: Preview/rule output is data. The skill instructs the agent to treat embedded text as untrusted (T7).
3. **Host agent → Subscription LLM**: Out of scope for `ocr-delegate`; user configures the agent's LLM provider.

### Threat Summary (aligned with SECURITY.md)

| ID | Threat | Mitigation |
|----|--------|------------|
| T1 | npm wrapper / prebuilt download chain exfiltration | Removed: no `npm i -g`, no `install.sh` prebuilt fetch, no auto-updater. Source build only (`make build`). |
| T2 | Telemetry export | `internal/telemetry` excluded from `ocr-delegate` dependency closure (CI Gate 3). |
| T3 | LLM provider endpoint strings in binary | `internal/llm` excluded from closure; CI Gate 4 string-scans artifact for forbidden provider domains. |
| T4 | Local viewer / DNS rebinding | `internal/viewer` excluded from closure; no HTTP listeners in delivered binary. |
| T5 | Session JSONL persistence residue | `internal/session` excluded; E2E verifies no new `~/.opencodereview` files after delegate runs. |
| T6 | FileReader absolute-path defect (upstream) | Upstream code unchanged per fork policy; `file_read` tool not in delegate closure. |
| T7 | Prompt injection via diff/rule content | Skill boundary: treat all preview output as data, not instructions. Adversarial fixtures in `testdata/adversarial/`. |
| T8 | Upstream GitHub Action (LLM exfiltration in CI) | `action.yml` removed; examples marked upstream-only. |
| T9 | Rule tampering (`.opencodereview/rule.json`) | Skill workflow flags rule changes as high-priority review targets. |

## Secure Design Principles

| Principle | How Applied (Delegate Edition) |
|-----------|--------------------------------|
| **Least privilege** | `CGO_ENABLED=0`. CLI requires no elevated permissions. Zero network capability. |
| **Fail-safe defaults** | No credentials in binary. No listeners. Source-first distribution. |
| **Complete mediation** | Git subprocesses use explicit argument lists; `--end-of-options` prevents flag injection. |
| **Economy of mechanism** | External execution limited to `git` with hardcoded subcommands — no shell invocation. |
| **Open design** | Apache-2.0. Security relies on verifiable closure and CI gates, not obscurity. |
| **Separation of privilege** | Engineering (`ocr-delegate`) separated from review intelligence (host agent LLM). |

## Countermeasures (Delegate Edition)

| Weakness | Applicability | Countermeasure |
|----------|---------------|----------------|
| **Command injection** (CWE-78) | `git` only, explicit args, no shell | Mitigated |
| **Path traversal** (CWE-22) | Delegate does not register upstream `file_read` tool (T6) | Not in delivery closure |
| **Credential leakage** (CWE-798) | No API keys in binary or config reads | Mitigated |
| **Vulnerable components** (CWE-1104) | `govulncheck` on delegate packages; Dependabot when enabled | Mitigated |
| **Prompt injection** | Skill discipline + adversarial regression fixtures (T7) | Mitigated at skill boundary |
| **Supply-chain / npm** (T1) | npm wrapper and prebuilt chain removed | Mitigated |

## Automated Verification

| Check | Tool | When |
|-------|------|------|
| Compile all packages | `go build ./...` | CI Gate 1 |
| Format | `gofmt -s -l` | CI (fmt-check) |
| Static analysis | `go vet` | CI / `make check` |
| Unit + race tests | `make test` (`LC_ALL=C`, extensions excluded) | CI Gate 2 |
| Coverage threshold | `make coverage` (≥ 80%) | CI |
| Known vulnerabilities | `govulncheck` on delegate closure | CI |
| Dependency closure | `go list -deps ./cmd/ocr-delegate` | CI Gate 3 |
| Provider string scan | `strings dist/ocr-delegate` | CI Gate 4 |
| Zero-network delegate | E2E `test/e2e/dryrun.sh` (S1, S4, S5, T5, adversarial) | CI |
| Skill/plugin stub integrity | CI plugin-skill check | CI |
| Dependency monitoring | Dependabot | Continuous (when enabled on fork) |

### CI Gate Mapping (`.github/workflows/ci.yml`)

| Gate | Step | Validates |
|------|------|-----------|
| **1** | `go build ./...` | All packages compile (upstream sync health) |
| **2** | `make test` + `make coverage` | Full test suite with race detector; coverage ≥ 80% |
| **3** | `go list -deps` closure check | No forbidden internal packages in `ocr-delegate` |
| **4** | `make build` + `strings` scan | No LLM provider domain strings in artifact |

Additional steps: `govulncheck`, `gofmt` check, `go vet`, E2E dry-run (S4: ripgrep for npm-free canonical install paths + delegate-only skill discipline).

---

## (Upstream Only) Extended Threat Model

The following applies to upstream **alibaba/open-code-review** (`opencodereview` / `ocr` CLI), not Delegate Edition.

### Additional Trust Boundaries (upstream)

- **CLI → LLM API**: API keys over HTTPS; responses validated (T2 upstream sense: key leakage).
- **CLI → Local output / viewer**: Path validation via `pathutil.WithinBase()`; host-header allowlist for viewer DNS rebinding (T4 upstream sense).

### Upstream Threat Notes

| Concern | Upstream mitigation |
|---------|---------------------|
| API key leakage | Keys from env only; never logged |
| Path traversal via LLM-suggested paths | `pathutil.WithinBase()` in file-read tool |
| DNS rebinding (viewer) | Host-header allowlist (`OCR_VIEWER_ALLOWED_HOSTS`) |
| Malicious LLM response | JSON schema validation; line bounds checking |
| TLS / MITM | Default Go TLS 1.2+; `InsecureSkipVerify` never set |

For the full upstream assurance narrative, see [alibaba/open-code-review](https://github.com/alibaba/open-code-review).
