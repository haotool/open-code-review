# CI/CD Integration Examples

> **Fork notice (haotool Delegate Edition):** Examples that install OCR via **npm** or call the upstream **`ocr`** CLI / GitHub Action document **alibaba/open-code-review** only. This fork ships **`ocr-delegate`** + the [agent skill](../skills/open-code-review/SKILL.md) — see the [README](../README.md#quick-install-3-steps). Upstream-only examples are marked below.

This directory contains examples for integrating OpenCodeReview (OCR) into various CI/CD pipelines.

## Contents

- **[github_actions/](./github_actions/)** — GitHub Actions integration (**upstream only**; `action.yml` removed in fork)
- **[gitlab_ci/](./gitlab_ci/)** — GitLab CI integration (**upstream only**; uses npm + `ocr review`)
- **[bitbucket_pipelines/](./bitbucket_pipelines/)** — Bitbucket Pipelines (**upstream only**; uses npm + `ocr review`)
- **[gitflic_ci/](./gitflic_ci/)** — GitFlic CI integration example (**upstream only**)
- **[gerrit_ci/](./gerrit_ci/)** — Gerrit (Jenkins / Gerrit Trigger) integration example (**upstream only**)

Each subdirectory contains its own README with detailed setup instructions.
