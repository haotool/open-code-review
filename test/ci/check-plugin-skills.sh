#!/usr/bin/env bash
# Verify plugin skills are complete and remain usable when the plugin is isolated.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"

canonical="${ROOT}/skills/open-code-review/SKILL.md"
delegate_canonical="${ROOT}/skills/open-code-review-delegate/SKILL.md"
marketplace="${ROOT}/.claude-plugin/marketplace.json"
plugin_main="${ROOT}/plugins/open-code-review/skills/open-code-review/SKILL.md"
plugin_delegate="${ROOT}/plugins/open-code-review/skills/open-code-review-delegate/SKILL.md"
claude_main="${ROOT}/plugins/open-code-review/claude-code/skills/open-code-review/SKILL.md"
claude_delegate="${ROOT}/plugins/open-code-review/claude-code/skills/open-code-review-delegate/SKILL.md"

for f in "${canonical}" "${delegate_canonical}" "${marketplace}" "${plugin_main}" "${plugin_delegate}" "${claude_main}" "${claude_delegate}"; do
  if [[ ! -f "${f}" ]]; then
    echo "FAIL: missing ${f}"
    exit 1
  fi
done

if ! python3 -m json.tool "${marketplace}" >/dev/null; then
  echo "FAIL: marketplace manifest is not valid JSON"
  exit 1
fi

for source in \
  './plugins/open-code-review' \
  './plugins/open-code-review/claude-code'; do
  if ! rg -q '"source"[[:space:]]*:[[:space:]]*"'"${source}"'"' "${marketplace}"; then
    echo "FAIL: marketplace manifest is missing plugin source ${source}"
    exit 1
  fi
done

if ! cmp -s "${canonical}" "${plugin_main}"; then
  echo "FAIL: plugin main skill must be a complete copy of the canonical skill"
  exit 1
fi

if ! cmp -s "${delegate_canonical}" "${plugin_delegate}"; then
  echo "FAIL: plugin delegate skill must be a complete copy of the canonical alias"
  exit 1
fi

if ! cmp -s "${canonical}" "${claude_main}"; then
  echo "FAIL: Claude plugin main skill must be a complete copy of the canonical skill"
  exit 1
fi

if ! cmp -s "${delegate_canonical}" "${claude_delegate}"; then
  echo "FAIL: Claude plugin delegate skill must be a complete copy of the canonical alias"
  exit 1
fi

for manifest in \
  "${ROOT}/plugins/open-code-review/.codex-plugin/plugin.json" \
  "${ROOT}/plugins/open-code-review/.cursor-plugin/plugin.json" \
  "${ROOT}/plugins/open-code-review/claude-code/.claude-plugin/plugin.json"; do
  test -f "${manifest}" || {
    echo "FAIL: missing plugin manifest ${manifest}"
    exit 1
  }
done

isolated="$(mktemp -d)"
trap 'rm -rf "${isolated}"' EXIT
cp -R "${ROOT}/plugins/open-code-review" "${isolated}/open-code-review"
package_root="${isolated}/open-code-review"

for f in \
  "${package_root}/skills/open-code-review/SKILL.md" \
  "${package_root}/skills/open-code-review-delegate/SKILL.md"; do
  test -s "${f}" || {
    echo "FAIL: isolated plugin package is missing ${f}"
    exit 1
  }
done

if rg -q '\.\./\.\./\.\./\.\./skills|repo-root' "${package_root}/skills"; then
  echo "FAIL: isolated plugin skill still depends on an external repo-root path"
  exit 1
fi

cmp -s "${canonical}" "${package_root}/skills/open-code-review/SKILL.md" || {
  echo "FAIL: isolated plugin main skill differs from canonical"
  exit 1
}
cmp -s "${delegate_canonical}" "${package_root}/skills/open-code-review-delegate/SKILL.md" || {
  echo "FAIL: isolated plugin delegate skill differs from canonical"
  exit 1
}

claude_package="${package_root}/claude-code"
for f in \
  "${claude_package}/.claude-plugin/plugin.json" \
  "${claude_package}/commands/delegate-review.md" \
  "${claude_package}/skills/open-code-review/SKILL.md" \
  "${claude_package}/skills/open-code-review-delegate/SKILL.md"; do
  test -s "${f}" || {
    echo "FAIL: isolated Claude plugin package is missing ${f}"
    exit 1
  }
done

if rg -q '\.\./\.\./\.\./\.\./skills|repo-root' "${claude_package}/skills"; then
  echo "FAIL: isolated Claude plugin skill still depends on an external repo-root path"
  exit 1
fi

cmp -s "${canonical}" "${claude_package}/skills/open-code-review/SKILL.md" || {
  echo "FAIL: isolated Claude plugin main skill differs from canonical"
  exit 1
}
cmp -s "${delegate_canonical}" "${claude_package}/skills/open-code-review-delegate/SKILL.md" || {
  echo "FAIL: isolated Claude plugin delegate skill differs from canonical"
  exit 1
}

echo "PASS: plugin skills complete and isolated-package safe"
