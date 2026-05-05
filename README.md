# ai-search-bot

Telegram-бот на Go. Сейчас работает как эхо-сервис — возвращает пользователю его же сообщение в виде reply.

## Требования

- Go 1.22+
- Telegram-бот, созданный через [@BotFather](https://t.me/BotFather)

## Быстрый старт

```bash
# Установить зависимости
go mod download

# Прописать токен в конфиг
# configs/config.dev.yaml → telegram.token

# Запустить
go run ./cmd/main.go
```

## Конфигурация

Конфиги лежат в `configs/`. Нужный файл выбирается по переменной окружения `APP_ENV` (по умолчанию `dev`):

| `APP_ENV` | Файл конфига             |
|-----------|--------------------------|
| `dev`     | `configs/config.dev.yaml`  |
| `prod`    | `configs/config.prod.yaml` |

Структура конфига:

```yaml
telegram:
  token: "ВАШ_ТОКЕН"   # токен от @BotFather

server:
  debug: false           # если true — бот логирует все HTTP-запросы к Telegram API

log:
  level: "info"          # debug | info | warn | error
  file: "logs/app.log"   # путь к лог-файлу; если не указан — логи идут в консоль
```

> Конфиг-файлы с токенами не коммитятся (прописаны в `.gitignore`).

## Среды

**dev** — логи текстом в консоль, `log.file` не нужен:

```yaml
log:
  level: "debug"
```

**prod** — логи JSON в файл, `log.file` обязателен:

```yaml
log:
  level: "info"
  file: "logs/app.log"
```

Запуск в prod:

```bash
APP_ENV=prod go run ./cmd/main.go
```

## Структура проекта

```
cmd/
  main.go                — точка входа: загрузка конфига, инициализация логгера, запуск бота
configs/
  config.dev.yaml
  config.prod.yaml
internal/
  config/
    config.go            — загрузка YAML-конфига по APP_ENV
  logger/
    logger.go            — фабрика: выбор реализации логгера по среде
    console.go           — текстовый логгер (stdout)
    file.go              — JSON-логгер в файл
  bot/
    bot.go               — обработка обновлений от Telegram, эхо-логика
```

## Остановка

`Ctrl+C` или `SIGTERM` — бот завершается gracefully: дожидается окончания текущей обработки, закрывает лог-файл.
