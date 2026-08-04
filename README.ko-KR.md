<div align="center">
  <h1>Open Code Review — Delegate Edition</h1>
  <p><strong><a href="https://github.com/haotool">haotool</a>의 보안 강화 Fork</strong></p>
</div>

<p align="center">
  <a href="README.md">English</a> | <a href="README.zh-CN.md">简体中文</a> | <a href="README.zh-TW.md">繁體中文</a> | <a href="README.ja-JP.md">日本語</a> | 한국어 | <a href="README.ru-RU.md">Русский</a>
</p>

<p align="center">
  <a href="https://github.com/haotool/open-code-review-delegate/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/haotool/open-code-review-delegate/ci.yml?style=flat-square" /></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/haotool/open-code-review-delegate?style=flat-square" /></a>
  <a href="https://github.com/alibaba/open-code-review"><img alt="Upstream" src="https://img.shields.io/badge/upstream-alibaba%2Fopen--code--review-blue?style=flat-square" /></a>
</p>

---

## 중요 안내

본 저장소는 Alibaba 공식 제품이 **아닙니다**. [haotool](https://github.com/haotool)이 Apache License 2.0 하에 독립 유지하는 Fork이며, [alibaba/open-code-review](https://github.com/alibaba/open-code-review)에서 파생되었습니다.

- 저장소: [`haotool/open-code-review-delegate`](https://github.com/haotool/open-code-review-delegate) (이전 `haotool/open-code-review` URL은 자동으로 리디렉션됩니다).
- Alibaba 상표를 사용하지 않으며 공식 보증을 암시하지 않습니다.
- upstream 저작권 표기는 [NOTICE](NOTICE)를 참조하세요.
- 전체 라이선스 조항은 [LICENSE](LICENSE)를 참조하세요.

## 개요

Open Code Review (Delegate Edition)는 단일 CLI 바이너리 **`ocr-delegate`**를 제공합니다. 결정론적 코드 리뷰 엔지니어링을 담당합니다.

- **파일 선별** — 변경된 파일 중 리뷰 대상을 정확히 결정
- **규칙 해석** — `rule.json`에서 파일별 리뷰 규칙 매칭
- **제로 네트워크** — 바이너리는 LLM 호출 및 아웃바운드 연결 없음
- **제로 자격 증명** — API 키, 토큰, 프로바이더 설정 불필요

AI 코딩 에이전트(Cursor, Claude Code, Codex 등)가 구독 LLM으로 실제 리뷰를 수행합니다. `ocr-delegate`는 엔지니어링, Host Agent는 지능을 담당합니다.

upstream 배포 모델의 npm 자동 업데이트, 사전 빌드 바이너리 다운로드, 내장 프로바이더 엔드포인트로 인한 공급망 위험을 제거합니다.

## 빠른 설치 (3단계)

소스에서 빌드 — npm, 사전 빌드 바이너리, API 키 불필요:

```bash
# 1. Clone
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate

# 2. ocr-delegate 빌드
make build

# 3. 에이전트 skill 설치
make install-skill
# Claude Code 사용자:
make install-skill SKILL_DIR=~/.claude/skills
```

검증:

```bash
which ocr-delegate && ocr-delegate -h
```

설치는 `test/e2e/dryrun.sh`로 검증됩니다 ([E2E 검증](test/e2e/README.md)).

## 사용법

`skills/open-code-review/SKILL.md`가 위임 워크플로의 단일 진실 공급원입니다.

### 1. Preview — 범위 결정

```bash
ocr-delegate preview [--from main --to feature] [--commit <hash>] [-b "context"]
```

### 2. 파일 규칙 조회

```bash
ocr-delegate rule <path1> <path2> ...
```

### 3. diff 및 리뷰

preview 출력의 mode/ref 메타데이터에 따라 git 사용. Host Agent가 구조화된 findings를 생성합니다.

전체 워크플로, 보안 규율(T7/T9), 출력 스키마는 [skill 문서](skills/open-code-review/SKILL.md)를 참조하세요.

공식 지원 Codex, Claude Code, Cursor 설치 경로는 [에이전트 지원 매트릭스](docs/AGENT_SUPPORT.md)를 참조하세요.

## 보안 posture

| 속성 | 보장 |
|------|------|
| 제로 아웃바운드 | CI 문자열 스캔 + E2E로 검증 |
| 제로 자격 증명 | Delegate 모드에서 LLM 프로바이더 설정 불필요 |
| 최소 의존성 클로저 | 금지 모듈(`internal/llm`, `telemetry`, `mcp`, `viewer`, `session`) 제외 |
| 소스 우선 배포 | npm 래퍼, 사전 빌드 다운로드, 자동 업데이터 없음 |
| 공급망 제거 | upstream 설치 스크립트, npm 패키지, GitHub Action 제거 |

자세히: [SECURITY.md](SECURITY.md)

## 빌드 재현 및 검증

본 Fork는 **소스 우선**입니다. 사전 빌드 바이너리를 배포하지 않습니다.

```bash
make build                          # dist/ocr-delegate 생성
shasum -a 256 dist/ocr-delegate     # 로컬 체크섬

# 크로스 플랫폼 로컬 빌드 (선택)
make dist                           # 모든 플랫폼 + sha256sum.txt
```

## 개발

```bash
make build       # ocr-delegate 빌드
make test        # 테스트 스위트 실행
make check       # fmt + vet + mod tidy
make coverage    # 커버리지 리포트 (80% 임계값)
```

[CONTRIBUTING.ko-KR.md](CONTRIBUTING.ko-KR.md)를 참조하세요.

## upstream 레거시 디렉터리

upstream에서 유지되었으나 본 Fork에서 **유지보수하지 않는** 디렉터리:

| 디렉터리 | 상태 |
|----------|------|
| `pages/` | upstream 랜딩 페이지 — 본 Fork에서 빌드/배포하지 않음 |
| `extensions/` | upstream VS Code 확장 — `ocr-delegate`가 아닌 전체 `ocr` CLI 필요 |

## 라이선스 및 귀속

[Apache-2.0](LICENSE). upstream Alibaba 기여자 저작권 표기는 [NOTICE](NOTICE)에 보존됩니다.

본 Fork는 upstream을 상당히 수정합니다. v1.0.0 변경 사항은 [CHANGELOG.md](CHANGELOG.md)를 참조하세요.

## 릴리스 체크리스트 (S1–S7)

| ID | 요구 | 검증 |
|----|------|------|
| S1 | 제로 아웃바운드 | `test/e2e/dryrun.sh` + CI Gate 4 문자열 스캔 |
| S2 | 제로 자격 | E2E dry-run + delegate 전용 바이너리 |
| S3 | 결정론 엔진 | 26 시나리오 golden 동등성 스위트 |
| S4 | 소스 우선 + delegate-only skill | 설치 경로에 npm/postinstall 공급망 없음; Skill은 `ocr-delegate`만 참조(전체 `ocr` CLI / 외부 LLM API 아님) |
| S5 | 설치 ≤3단계 | E2E dry-run 3단계 검증 |
| S6 | 테스트 통과 | CI Gate 2 (`go test ./...`) |
| S7 | 주입 저항 | 적대 fixture 라이브러리 + E2E spot-check |
