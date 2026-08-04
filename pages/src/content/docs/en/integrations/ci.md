---
title: CI/CD
sidebar:
  order: 4
---

> **Delegate Edition (this fork):** Automated CI pipelines that install the upstream `ocr` CLI are **not** the canonical workflow here. Build **`ocr-delegate`** from source and run reviews via the agent skill locally — see [Delegation Mode](../delegate/).

## Delegate Edition — CI/CD

This fork does **not** ship npm-based CI recipes. `ocr-delegate` is a zero-network binary for deterministic file selection and rule resolution; the host agent performs LLM reasoning with its subscription quota.

### Recommended approach

Run code review in your development flow or agent session:

```bash
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate
make build
export PATH="$PWD/dist:$PATH"
make install-skill
```

In the repository under review:

```bash
ocr-delegate preview --from origin/main --to HEAD
ocr-delegate rule <paths-from-preview>
```

Use git for diffs based on preview mode/ref metadata. Your host agent (Cursor, Claude Code, Codex, etc.) completes the review — see the [skill workflow](../delegate/).

### Optional: self-hosted runner with source build

If you need automation without npm, build `ocr-delegate` on the runner and invoke `preview` / `rule` in a job script. Do **not** expect JSON review comments from `ocr-delegate` alone — post-processing and LLM review remain the host agent's responsibility.

Example sketch (GitHub Actions):

```yaml
- name: Build ocr-delegate
  run: make build
- name: Preview changed files
  run: ./dist/ocr-delegate preview --from origin/${{ github.base_ref }} --to HEAD
```

Extend with your agent pipeline as needed.

## See Also

- [Delegation Mode](../delegate/) — preview → rule → review workflow.
- [QuickStart](../../quickstart/) — three-step install.
- [E2E dry-run](https://github.com/haotool/open-code-review-delegate/tree/main/test/e2e) — zero-outbound verification.

---

## Upstream CI/CD (not shipped in this fork)

> The following applies to **upstream [alibaba/open-code-review](https://github.com/alibaba/open-code-review)** only.

Upstream provides ready-made GitHub Actions and GitLab CI pipelines that install the `ocr` CLI, configure LLM credentials from secrets, run `ocr review --format json`, and post inline PR/MR comments.

- GitHub Actions: [`examples/github_actions/ocr-review.yml`](https://github.com/alibaba/open-code-review/blob/main/examples/github_actions/ocr-review.yml)
- GitLab CI: [`examples/gitlab_ci/.gitlab-ci.yml`](https://github.com/alibaba/open-code-review/blob/main/examples/gitlab_ci/.gitlab-ci.yml)

See the upstream documentation for secrets, customization, and troubleshooting.
