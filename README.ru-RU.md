<div align="center">
  <h1>Open Code Review — Delegate Edition</h1>
  <p><strong>Форк с усиленной безопасностью от <a href="https://github.com/haotool">haotool</a></strong></p>
</div>

<p align="center">
  <a href="README.md">English</a> | <a href="README.zh-CN.md">简体中文</a> | <a href="README.zh-TW.md">繁體中文</a> | <a href="README.ja-JP.md">日本語</a> | <a href="README.ko-KR.md">한국어</a> | Русский
</p>

<p align="center">
  <a href="https://github.com/haotool/open-code-review-delegate/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/haotool/open-code-review-delegate/ci.yml?style=flat-square" /></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/haotool/open-code-review-delegate?style=flat-square" /></a>
  <a href="https://github.com/alibaba/open-code-review"><img alt="Upstream" src="https://img.shields.io/badge/upstream-alibaba%2Fopen--code--review-blue?style=flat-square" /></a>
</p>

---

## Важное уведомление

Этот репозиторий **не** является официальным продуктом Alibaba. Это независимый форк [haotool](https://github.com/haotool), производный от [alibaba/open-code-review](https://github.com/alibaba/open-code-review) под Apache License 2.0.

- Репозиторий: [`haotool/open-code-review-delegate`](https://github.com/haotool/open-code-review-delegate) (запросы к прежнему URL `haotool/open-code-review` перенаправляются автоматически).
- Мы не используем товарные знаки Alibaba и не подразумеваем официальную поддержку.
- Атрибуция авторских прав upstream — в [NOTICE](NOTICE).
- Полные условия лицензии — в [LICENSE](LICENSE).

## Что это?

Open Code Review (Delegate Edition) поставляет один CLI-бинарник **`ocr-delegate`** для детерминированной инженерии code review:

- **Отбор файлов** — точное определение изменённых файлов для ревью
- **Разрешение правил** — сопоставление правил из `rule.json` с каждым файлом
- **Нулевой исходящий трафик** — бинарник не вызывает LLM и не открывает исходящие соединения
- **Нулевые учётные данные** — API-ключи, токены и конфигурация провайдера не требуются

AI-агент (Cursor, Claude Code, Codex и т.д.) выполняет реальное ревью через свой LLM по подписке. `ocr-delegate` отвечает за инженерию; host agent — за интеллект.

Такая модель устраняет риски цепочки поставок от npm auto-updater, скачивания prebuilt-бинарников и встроенных endpoint провайдера в upstream-модели распространения.

## Быстрая установка (3 шага)

Сборка из исходников — без npm, prebuilt-бинарников и API-ключей:

```bash
# 1. Clone
git clone https://github.com/haotool/open-code-review-delegate.git && cd open-code-review-delegate

# 2. Сборка ocr-delegate
make build

# 3. Установка agent skill
make install-skill
# Для Claude Code:
make install-skill SKILL_DIR=~/.claude/skills
```

Проверка:

```bash
which ocr-delegate && ocr-delegate -h
```

Установка проверяется `test/e2e/dryrun.sh` ([E2E-верификация](test/e2e/README.md)).

## Использование

`skills/open-code-review/SKILL.md` — единый источник истины для delegate workflow.

### 1. Preview — определить область

```bash
ocr-delegate preview [--from main --to feature] [--commit <hash>] [-b "context"]
```

### 2. Правила для файлов

```bash
ocr-delegate rule <path1> <path2> ...
```

### 3. Diff и ревью

Используйте git по mode/ref из preview. Host agent формирует структурированные findings.

Полный workflow, дисциплина безопасности (T7/T9) и схема вывода — в [документации skill](skills/open-code-review/SKILL.md).

Формально поддерживаемые пути установки для Codex, Claude Code и Cursor — в [матрице поддержки агентов](docs/AGENT_SUPPORT.md).

## Позиция по безопасности

| Свойство | Гарантия |
|----------|----------|
| Нулевой исходящий | Без сетевых соединений (CI string scan + E2E) |
| Нулевые учётные данные | Конфигурация LLM-провайдера не нужна в delegate mode |
| Минимальное замыкание зависимостей | Запрещённые модули (`internal/llm`, `telemetry`, `mcp`, `viewer`, `session`) исключены |
| Распространение из исходников | Без npm-обёртки, prebuilt-скачиваний и auto-updater |
| Обрезанная цепочка поставок | Удалены upstream install scripts, npm-пакеты и GitHub Action |

Подробнее: [SECURITY.md](SECURITY.md)

## Воспроизведение сборки

Этот форк **source-first**: prebuilt-бинарники не публикуются.

```bash
make build                          # создаёт dist/ocr-delegate
shasum -a 256 dist/ocr-delegate     # локальная контрольная сумма

# Кросс-платформенная локальная сборка (опционально)
make dist                           # все платформы + sha256sum.txt
```

## Разработка

```bash
make build       # сборка ocr-delegate
make test        # тестовый набор
make check       # fmt + vet + mod tidy
make coverage    # отчёт покрытия (порог 80%)
```

См. [CONTRIBUTING.ru-RU.md](CONTRIBUTING.ru-RU.md).

## Legacy-директории upstream

Сохранены из upstream, но **не поддерживаются** в этом форке:

| Каталог | Статус |
|---------|--------|
| `pages/` | Landing upstream — не собирается и не деплоится этим форком |
| `extensions/` | VS Code extension upstream — требует полный `ocr` CLI, не `ocr-delegate` |

## Лицензия и атрибуция

[Apache-2.0](LICENSE). Уведомления об авторских правах upstream Alibaba сохранены в [NOTICE](NOTICE).

Форк существенно изменяет upstream. Изменения v1.0.0 — в [CHANGELOG.md](CHANGELOG.md).

## Чеклист релиза (S1–S7)

| ID | Требование | Проверка |
|----|------------|----------|
| S1 | Нулевая исходящая сеть | `test/e2e/dryrun.sh` + CI Gate 4 string scan |
| S2 | Нулевые ключи | E2E dry-run + delegate-only binary |
| S3 | Детерминированный движок | Golden suite из 26 сценариев |
| S4 | Source-first + delegate-only skill | Без npm/postinstall в путях установки; skill ссылается только на `ocr-delegate` (не полный `ocr` CLI / внешний LLM API) |
| S5 | Установка ≤3 шагов | E2E dry-run трёхшаговая проверка |
| S6 | Тесты проходят | CI Gate 2 (`go test ./...`) |
| S7 | Устойчивость к инъекциям | Adversarial fixture library + E2E spot-checks |
