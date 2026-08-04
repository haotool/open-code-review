# Участие в разработке — Open Code Review (Delegate Edition)

Спасибо за интерес! Форк поддерживается [haotool](https://github.com/haotool), производный от [alibaba/open-code-review](https://github.com/alibaba/open-code-review) (Apache-2.0).

## Настройка

```bash
git clone https://github.com/haotool/open-code-review-delegate.git
cd open-code-review-delegate
make build && make test
```

## Команды разработки

```bash
make build && make test && make check && make coverage
make install-skill SKILL_DIR=~/.cursor/skills
```

## Правила

- Следуйте существующим Go-соглашениям
- Перед PR: `make check`
- Тесты для нового поведения
- Сохраняйте замыкание зависимостей `ocr-delegate` — не импортируйте `internal/llm`, `telemetry`, `mcp`, `viewer`, `session` из delegate-путей
- Conventional Commits

## Область

- **`ocr-delegate`** — основной артеfact; не расширяйте сеть/учётные данные
- **`pages/`**, **`extensions/`** — legacy upstream, не поддерживаются в форке
- Не возвращайте npm, prebuilt downloads, upstream LLM Action

## Лицензия

Вклад лицензируется под [Apache-2.0](LICENSE).
