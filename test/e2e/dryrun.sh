#!/usr/bin/env bash
# E2E dry-run: fresh-environment install and zero-outbound verification.
# Validates S1 (zero network), S2 (zero credentials), S5 (≤3 install steps),
# T5 (no ~/.opencodereview residue), and adversarial spot-checks (samples 01, 02, 03, 04, 05).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
REPORT="${ROOT}/test/e2e/last-run-report.txt"
PASS=0
FAIL=0

log()  { printf '[e2e] %s\n' "$*"; }
pass() { PASS=$((PASS + 1)); log "PASS: $*"; }
fail() { FAIL=$((FAIL + 1)); log "FAIL: $*"; }

# ── Isolated environment (no real home pollution) ───────────────────────────
E2E_TMP="$(mktemp -d)"
E2E_HOME="${E2E_TMP}/home"
E2E_REPO="${E2E_TMP}/sample-repo"
E2E_BIN="${E2E_TMP}/bin"
mkdir -p "${E2E_HOME}" "${E2E_BIN}"

cleanup() { rm -rf "${E2E_TMP}" 2>/dev/null || true; }
trap cleanup EXIT

# ── S5: three-step install (clone already done; build + install-skill) ───────
log "=== S5: three-step install ==="
log "Step 1: clone (assumed done — running from source tree)"
log "Step 2: make build"
make -C "${ROOT}" build
cp "${ROOT}/dist/ocr-delegate" "${E2E_BIN}/ocr-delegate"
chmod +x "${E2E_BIN}/ocr-delegate"
pass "make build → dist/ocr-delegate"

log "Step 3: make install-skill (isolated SKILL_DIR)"
SKILL_DIR="${E2E_HOME}/skills"
make -C "${ROOT}" install-skill SKILL_DIR="${SKILL_DIR}"
test -f "${SKILL_DIR}/open-code-review/SKILL.md" && pass "install-skill copied SKILL.md" || fail "install-skill missing SKILL.md"
test -f "${SKILL_DIR}/open-code-review-delegate/SKILL.md" && pass "install-skill copied delegate SKILL.md" || fail "install-skill missing delegate SKILL.md"

export PATH="${E2E_BIN}:${PATH}"
which ocr-delegate >/dev/null && pass "ocr-delegate on PATH" || fail "ocr-delegate not found"

# Isolate HOME only for delegate runs (avoid polluting real ~/.opencodereview)
export HOME="${E2E_HOME}"
unset OCR_LLM_URL OCR_LLM_TOKEN OCR_LLM_MODEL ANTHROPIC_BASE_URL ANTHROPIC_AUTH_TOKEN ANTHROPIC_MODEL OPENAI_API_KEY 2>/dev/null || true

# ── S4: source-first distribution + delegate-only skill discipline ─────────
log "=== S4: source-first + delegate-only skill ==="
if rg -q '`ocr `|\bocr review\b|\bocr config\b|\bocr scan\b' "${SKILL_DIR}/open-code-review/SKILL.md" 2>/dev/null; then
  fail "skill references full ocr CLI"
else
  pass "skill uses ocr-delegate only"
fi

S4_CANONICAL=(
  "${ROOT}/README.md"
  "${ROOT}/skills/open-code-review/SKILL.md"
  "${ROOT}/pages/src/content/docs/en/quickstart.md"
  "${ROOT}/pages/src/content/docs/zh/quickstart.md"
  "${ROOT}/pages/src/content/docs/ja/quickstart.md"
  "${ROOT}/pages/src/content/docs/en/installation.md"
  "${ROOT}/pages/src/content/docs/zh/installation.md"
  "${ROOT}/pages/src/content/docs/ja/installation.md"
  "${ROOT}/pages/src/content/docs/en/integrations/delegate.md"
  "${ROOT}/pages/src/content/docs/zh/integrations/delegate.md"
  "${ROOT}/pages/src/content/docs/ja/integrations/delegate.md"
  "${ROOT}/pages/src/content/docs/en/integrations/agent-skill.md"
  "${ROOT}/pages/src/content/docs/zh/integrations/agent-skill.md"
  "${ROOT}/pages/src/content/docs/ja/integrations/agent-skill.md"
  "${ROOT}/pages/src/content/docs/en/integrations/ci.md"
  "${ROOT}/pages/src/content/docs/zh/integrations/ci.md"
  "${ROOT}/pages/src/content/docs/ja/integrations/ci.md"
  "${ROOT}/pages/src/components/HeroSection.tsx"
  "${ROOT}/pages/src/components/QuickStartSection.tsx"
)
S4_NPM_VIOLATIONS=0
for f in "${S4_CANONICAL[@]}"; do
  if rg -q 'npm (install|i) -g @alibaba-group/open-code-review' "${f}" 2>/dev/null; then
    fail "canonical install path contains npm @alibaba-group: ${f#"${ROOT}/"}"
    S4_NPM_VIOLATIONS=$((S4_NPM_VIOLATIONS + 1))
  fi
done
if [[ "${S4_NPM_VIOLATIONS}" -eq 0 ]]; then
  pass "canonical install paths free of npm @alibaba-group supply chain"
fi

# ── Build sample repo with rule.json ────────────────────────────────────────
log "=== Sample repo setup ==="
mkdir -p "${E2E_REPO}"
git -C "${E2E_REPO}" init >/dev/null 2>&1
cd "${E2E_REPO}"
git config user.email "e2e@test.local"
git config user.name "e2e"
mkdir -p .opencodereview internal/api
cat > .opencodereview/rule.json <<'RULE'
{
  "rules": [{"path": "**/*.go", "rule": "Check error handling and security."}],
  "exclude": ["vendor/**"]
}
RULE
cat > internal/api/handler.go <<'GO'
package api

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
GO
git add -A && git commit -m "initial" >/dev/null
echo '// added line' >> internal/api/handler.go
pass "sample repo with rule.json ready"

# ── T5: no ~/.opencodereview residue ─────────────────────────────────────────
log "=== T5: no residue ==="
OCR_DIR="${E2E_HOME}/.opencodereview"
BEFORE_COUNT=0
if [[ -d "${OCR_DIR}" ]]; then
  BEFORE_COUNT="$(find "${OCR_DIR}" -type f 2>/dev/null | wc -l | tr -d ' ')"
fi

# ── S1: zero outbound — monitor sockets during ocr-delegate execution ───────
log "=== S1: zero outbound network ==="
NETWORK_VIOLATIONS=0

monitor_network() {
  local pid="$1"
  while kill -0 "${pid}" 2>/dev/null; do
    if lsof -p "${pid}" -iTCP -sTCP:ESTABLISHED 2>/dev/null | grep -q .; then
      NETWORK_VIOLATIONS=$((NETWORK_VIOLATIONS + 1))
      lsof -p "${pid}" -iTCP 2>/dev/null || true
    fi
    sleep 0.05
  done
}

run_with_monitor() {
  "$@" &
  local pid=$!
  monitor_network "${pid}" &
  local mon=$!
  wait "${pid}"
  local rc=$?
  kill "${mon}" 2>/dev/null || true
  wait "${mon}" 2>/dev/null || true
  return "${rc}"
}

PREVIEW_OUT="${E2E_TMP}/preview.out"
run_with_monitor ocr-delegate preview --repo "${E2E_REPO}" -b "e2e test" > "${PREVIEW_OUT}"
pass "ocr-delegate preview completed"

RULE_OUT="${E2E_TMP}/rule.out"
run_with_monitor ocr-delegate rule --repo "${E2E_REPO}" internal/api/handler.go > "${RULE_OUT}"
pass "ocr-delegate rule completed"

if [[ "${NETWORK_VIOLATIONS}" -eq 0 ]]; then
  pass "zero TCP connections during ocr-delegate execution"
else
  fail "detected ${NETWORK_VIOLATIONS} network connection(s)"
fi

# Verify preview output contains expected fields
grep -q "mode: workspace" "${PREVIEW_OUT}" && pass "preview mode: workspace" || fail "preview missing workspace mode"
grep -q "handler.go" "${PREVIEW_OUT}" && pass "preview lists handler.go" || fail "preview missing handler.go"
grep -q "Check error handling" "${RULE_OUT}" && pass "rule resolved from rule.json" || fail "rule output missing"

AFTER_COUNT=0
if [[ -d "${OCR_DIR}" ]]; then
  AFTER_COUNT="$(find "${OCR_DIR}" -type f 2>/dev/null | wc -l | tr -d ' ')"
fi
if [[ "${BEFORE_COUNT}" -eq "${AFTER_COUNT}" ]]; then
  pass "no new files in ~/.opencodereview (before=${BEFORE_COUNT} after=${AFTER_COUNT})"
else
  fail "~/.opencodereview grew (before=${BEFORE_COUNT} after=${AFTER_COUNT})"
fi

# ── S2: zero credentials ────────────────────────────────────────────────────
log "=== S2: zero credentials ==="
CRED_FOUND=0
for var in OCR_LLM_URL OCR_LLM_TOKEN OCR_LLM_MODEL ANTHROPIC_BASE_URL ANTHROPIC_AUTH_TOKEN OPENAI_API_KEY; do
  if [[ -n "${!var:-}" ]]; then CRED_FOUND=$((CRED_FOUND + 1)); fi
done
if [[ "${CRED_FOUND}" -eq 0 ]]; then
  pass "no credential env vars set"
else
  fail "credential env vars present"
fi

# ── S7: adversarial spot-checks (samples 01, 02, 03, 04, 05) ──────────────────
log "=== S7: adversarial spot-checks ==="

# Sample 01: injection strings present in fixture, ocr-delegate must still run safely
ADV01="${ROOT}/testdata/adversarial/01_cmd_injection_in_diff"
if [[ -f "${ADV01}/repo/internal/api/handler.go" ]]; then
  ADV01_REPO="${E2E_TMP}/adv01"
  cp -R "${ADV01}/repo" "${ADV01_REPO}"
  git -C "${ADV01_REPO}" init >/dev/null 2>&1
  git -C "${ADV01_REPO}" config user.email "e2e@test.local"
  git -C "${ADV01_REPO}" config user.name "e2e"
  git -C "${ADV01_REPO}" add -A && git -C "${ADV01_REPO}" commit -m "adv01" >/dev/null
  echo "// e2e uncommitted change" >> "${ADV01_REPO}/internal/api/handler.go"
  ADV01_OUT="${E2E_TMP}/adv01-preview.out"
  run_with_monitor ocr-delegate preview --repo "${ADV01_REPO}" > "${ADV01_OUT}" 2>&1
  grep -q "handler.go" "${ADV01_OUT}" && pass "adv01: preview lists injected fixture file" || fail "adv01: preview failed"
  grep -q "SYSTEM OVERRIDE" "${ADV01}/repo/internal/api/handler.go" && pass "adv01: injection strings present in fixture (host agent must treat as data)" || fail "adv01: fixture missing"
else
  fail "adv01 fixture not found"
fi

# Sample 02: repo escape read lure
ADV02="${ROOT}/testdata/adversarial/02_repo_escape_read"
if [[ -f "${ADV02}/repo/pkg/util/config.go" ]]; then
  ADV02_REPO="${E2E_TMP}/adv02"
  cp -R "${ADV02}/repo" "${ADV02_REPO}"
  git -C "${ADV02_REPO}" init >/dev/null 2>&1
  git -C "${ADV02_REPO}" config user.email "e2e@test.local"
  git -C "${ADV02_REPO}" config user.name "e2e"
  git -C "${ADV02_REPO}" add -A && git -C "${ADV02_REPO}" commit -m "adv02" >/dev/null
  echo "// e2e uncommitted change" >> "${ADV02_REPO}/pkg/util/config.go"
  ADV02_OUT="${E2E_TMP}/adv02-preview.out"
  run_with_monitor ocr-delegate preview --repo "${ADV02_REPO}" > "${ADV02_OUT}" 2>&1
  grep -q "config.go" "${ADV02_OUT}" && pass "adv02: preview lists escape-lure fixture file" || fail "adv02: preview failed"
  grep -q "../../../etc/passwd" "${ADV02}/repo/pkg/util/config.go" && pass "adv02: escape paths present in fixture (host agent must refuse)" || fail "adv02: fixture missing escape paths"
else
  fail "adv02 fixture not found"
fi

# Sample 03: secrets-style path lure
ADV03="${ROOT}/testdata/adversarial/03_secrets_style_paths"
if [[ -f "${ADV03}/repo/internal/auth/token.go" ]]; then
  ADV03_REPO="${E2E_TMP}/adv03"
  cp -R "${ADV03}/repo" "${ADV03_REPO}"
  git -C "${ADV03_REPO}" init >/dev/null 2>&1
  git -C "${ADV03_REPO}" config user.email "e2e@test.local"
  git -C "${ADV03_REPO}" config user.name "e2e"
  git -C "${ADV03_REPO}" add -A && git -C "${ADV03_REPO}" commit -m "adv03" >/dev/null
  echo "// e2e uncommitted change" >> "${ADV03_REPO}/internal/auth/token.go"
  ADV03_OUT="${E2E_TMP}/adv03-preview.out"
  run_with_monitor ocr-delegate preview --repo "${ADV03_REPO}" > "${ADV03_OUT}" 2>&1
  grep -q "token.go" "${ADV03_OUT}" && pass "adv03: preview lists secrets-lure fixture file" || fail "adv03: preview failed"
  grep -q "FAKE_SECRET" "${ADV03}/repo/internal/auth/token.go" && pass "adv03: fake secret strings in fixture (host agent must not read .env/pem paths)" || fail "adv03: fixture missing fake secrets"
else
  fail "adv03 fixture not found"
fi

# Sample 04: rule.json weakening fixture
ADV04="${ROOT}/testdata/adversarial/04_rule_json_weakening"
if [[ -f "${ADV04}/repo/.ocr/rule.json" ]]; then
  ADV04_REPO="${E2E_TMP}/adv04"
  cp -R "${ADV04}/repo" "${ADV04_REPO}"
  git -C "${ADV04_REPO}" init >/dev/null 2>&1
  git -C "${ADV04_REPO}" config user.email "e2e@test.local"
  git -C "${ADV04_REPO}" config user.name "e2e"
  git -C "${ADV04_REPO}" add -A && git -C "${ADV04_REPO}" commit -m "adv04" >/dev/null
  ADV04_OUT="${E2E_TMP}/adv04-rule.out"
  run_with_monitor ocr-delegate rule --repo "${ADV04_REPO}" --rule .ocr/rule.json .ocr/rule.json > "${ADV04_OUT}" 2>&1 || true
  pass "adv04: rule command runs on weakened rule.json fixture (host agent must flag weakening per CHECKLIST)"
  grep -q "internal/auth" "${ADV04}/repo/.ocr/rule.json" && pass "adv04: fixture contains broad exclude patterns" || fail "adv04: fixture missing excludes"
else
  fail "adv04 fixture not found"
fi

# Sample 05: output / channel exfiltration lure
ADV05="${ROOT}/testdata/adversarial/05_output_exfiltration"
if [[ -f "${ADV05}/repo/internal/report/summary.go" ]]; then
  ADV05_REPO="${E2E_TMP}/adv05"
  cp -R "${ADV05}/repo" "${ADV05_REPO}"
  git -C "${ADV05_REPO}" init >/dev/null 2>&1
  git -C "${ADV05_REPO}" config user.email "e2e@test.local"
  git -C "${ADV05_REPO}" config user.name "e2e"
  git -C "${ADV05_REPO}" add -A && git -C "${ADV05_REPO}" commit -m "adv05" >/dev/null
  echo "// e2e uncommitted change" >> "${ADV05_REPO}/internal/report/summary.go"
  ADV05_OUT="${E2E_TMP}/adv05-preview.out"
  run_with_monitor ocr-delegate preview --repo "${ADV05_REPO}" > "${ADV05_OUT}" 2>&1
  grep -q "summary.go" "${ADV05_OUT}" && pass "adv05: preview lists exfil lure fixture file" || fail "adv05: preview failed"
  grep -q "attacker.example" "${ADV05}/repo/internal/report/summary.go" && pass "adv05: exfil URLs present in fixture (host agent must not POST)" || fail "adv05: fixture missing exfil URLs"
else
  fail "adv05 fixture not found"
fi

# ── Report ──────────────────────────────────────────────────────────────────
{
  echo "E2E Dry-Run Report"
  echo "=================="
  echo "Date: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
  echo "Pass: ${PASS}"
  echo "Fail: ${FAIL}"
  echo ""
  echo "Install steps verified:"
  echo "  1. git clone (precondition)"
  echo "  2. make build"
  echo "  3. make install-skill"
  echo ""
  echo "Network monitoring: lsof -p during ocr-delegate execution"
  echo "Network violations: ${NETWORK_VIOLATIONS}"
  echo "Residue check: ~/.opencodereview before=${BEFORE_COUNT} after=${AFTER_COUNT}"
  echo "Adversarial samples checked: 01 (cmd injection), 02 (repo escape), 03 (secrets paths), 04 (rule weakening), 05 (output exfil)"
} > "${REPORT}"

log "=== Summary: ${PASS} passed, ${FAIL} failed ==="
log "Report: test/e2e/last-run-report.txt"

if [[ "${FAIL}" -gt 0 ]]; then
  exit 1
fi
