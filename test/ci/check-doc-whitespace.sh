#!/usr/bin/env bash
# Fail if canonical docs/components contain trailing whitespace (style debt guard).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"

CANONICAL=(
  "${ROOT}/ASSURANCE_CASE.md"
  "${ROOT}/SECURITY.md"
  "${ROOT}/README.md"
  "${ROOT}/README.zh-TW.md"
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

violations=0
for f in "${CANONICAL[@]}"; do
  if rg -q '[ \t]+$' "${f}" 2>/dev/null; then
    echo "FAIL: trailing whitespace: ${f#"${ROOT}/"}"
    violations=$((violations + 1))
  fi
done

if [[ "${violations}" -gt 0 ]]; then
  echo "doc-whitespace-check: ${violations} file(s) with trailing whitespace"
  exit 1
fi

echo "doc-whitespace-check passed"
