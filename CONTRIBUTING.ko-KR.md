# 기여 가이드 — Open Code Review (Delegate Edition)

[haotool](https://github.com/haotool)이 [alibaba/open-code-review](https://github.com/alibaba/open-code-review)(Apache-2.0)에서 파생한 Fork에 기여해 주셔서 감사합니다.

## 설정

```bash
git clone https://github.com/haotool/open-code-review-delegate.git
cd open-code-review-delegate
make build && make test
```

## 개발 명령

```bash
make build && make test && make check && make coverage
make install-skill SKILL_DIR=~/.cursor/skills
```

## 가이드라인

- 기존 Go 규약 준수
- PR 전 `make check` 실행
- 새 동작에는 테스트 추가
- `ocr-delegate` 의존성 클로저 유지 — delegate 경로에서 `internal/llm`, `telemetry`, `mcp`, `viewer`, `session` import 금지
- Conventional Commits

## 범위

- **`ocr-delegate`** 가 주요 산출물; 네트워크/자격 증명 표면 확대 금지
- **`pages/`**, **`extensions/`** 는 upstream 레거시 — 본 Fork 비유지
- npm 배포, 프리빌트 다운로드, upstream LLM Action 재도입 금지

## 라이선스

기여는 [Apache-2.0](LICENSE) 하에 제공됩니다.
