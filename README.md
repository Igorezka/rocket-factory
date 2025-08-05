![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/Igorezka/63e8befdd35f9d8a520926772212e73d/raw/coverage.json)
# Rocket factory

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0](https://olezhek28.courses/microservices) от [Олега Козырева](http://t.me/olezhek28go).

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
