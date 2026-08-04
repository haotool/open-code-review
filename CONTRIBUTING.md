# Contributing to Open Code Review (Delegate Edition)

Thank you for your interest in contributing! This fork is maintained by [haotool](https://github.com/haotool) and derived from [alibaba/open-code-review](https://github.com/alibaba/open-code-review) under Apache-2.0.

## Code of Conduct

By participating, you agree to maintain a respectful and inclusive environment.

## Ways to Contribute

- **Report bugs** — Open an issue with reproduction steps
- **Suggest improvements** — Open a feature request issue
- **Improve documentation** — Fix typos, clarify explanations, add examples
- **Write code** — Fix bugs, add tests, improve the delegate CLI

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Git](https://git-scm.com/) >= 2.41
- [Make](https://www.gnu.org/software/make/)

### Setup

```bash
git clone https://github.com/haotool/open-code-review-delegate.git
cd open-code-review-delegate
make build
make test
```

### Development Commands

```bash
make build       # build ocr-delegate to dist/
make build-full  # build full opencodereview CLI (development only)
make test        # run test suite with race detector
make check       # fmt + vet + mod tidy
make coverage    # coverage report (80% threshold)
make install-skill SKILL_DIR=~/.cursor/skills  # install agent skill locally
```

## Coding Guidelines

- Follow existing Go conventions in the codebase
- Run `make check` before submitting
- Add tests for new behavior
- Keep `ocr-delegate` dependency closure clean — do not import `internal/llm`, `telemetry`, `mcp`, `viewer`, or `session` from delegate code paths
- Commit messages: conventional format (`feat(scope):`, `fix(scope):`, `docs:`, `test:`, `ci:`, `chore:`)

## Test quarantine (resolved)

`TestCommentWorkerPool_AwaitKeyConcurrentSubmitOtherKey` previously hung under
`go test -race` when producers used a busy loop. The test now uses bounded
producers and runs in the full suite (`make test` and CI) without `-skip`.
See issue #12 for history.

## Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make changes with tests
4. Run `make check && make test`
5. Submit a PR with a clear description

## Scope Notes

- **`ocr-delegate`** is the primary deliverable; changes should not expand its network or credential surface
- **`pages/`** and **`extensions/`** are upstream legacy directories — not actively maintained in this fork
- Do not reintroduce npm distribution, prebuilt binary downloads, or upstream GitHub Action workflows

## License

By contributing, you agree that your contributions will be licensed under the [Apache-2.0 License](LICENSE).
