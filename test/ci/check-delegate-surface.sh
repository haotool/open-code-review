#!/usr/bin/env bash
# Gate 5: prevent removed full-distribution entry points from returning.
set -euo pipefail

ROOT="${OCR_DELEGATE_REPO_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"

if [[ ! -d "${ROOT}" ]]; then
  echo "FAIL: repository root does not exist: ${ROOT}"
  exit 1
fi

found=()
while IFS= read -r -d '' file; do
  rel="${file#"${ROOT}"/}"
  case "${rel}" in
    action.yml|*/action.yml|bin/ocr.js|*/bin/ocr.js|ocr.js|*/ocr.js)
      found+=("${rel}")
      ;;
  esac
done < <(
  find "${ROOT}" \
    -path "${ROOT}/.git" -prune -o \
    -path "${ROOT}/.worktrees" -prune -o \
    -type f -print0
)

if ((${#found[@]} > 0)); then
  echo "FAIL: removed full-distribution entry points are present:"
  printf '  %s\n' "${found[@]}"
  echo "Delegate Edition must not reintroduce action.yml, bin/ocr.js, or ocr.js."
  exit 1
fi

echo "PASS: Gate 5 delegate surface is free of action.yml/bin/ocr.js/ocr.js"
