# Uptime Monitor

Небольшой сервис мониторинга доступности URL. Периодически опрашивает адреса, пишет историю в Postgres, отдает JSON API
и метрики для Prometheus/Grafana. В репо есть docker-compose с готовыми Prometheus/Grafana/Loki.

## Возможности

- Периодические проверки URL с настраиваемыми интервалом, таймаутом и параллелизмом.
- Добавление/удаление и активация/деактивация URL.
- История по каждому URL (статус-код и задержка).
- Метрики `/api/metrics` (Prometheus) и готовый стек наблюдаемости.
- Авто-применение миграций Goose при старте.

## Быстрый запуск через Docker

Требуется Docker + Docker Compose.

1. Подготовьте конфиг, который читает сервис (`config/config.prod.yaml`). Для старта можно скопировать dev-шаблон:
   ```sh
   cp config/config.dev.yaml config/config.prod.yaml
   ```
2. Запустите стек:
   ```sh
   docker-compose up --build
   ```
    - API: http://localhost:8080/api
    - Prometheus: http://localhost:9090
    - Grafana: http://localhost:3000 (логин/пароль admin/admin)
3. Добавьте URL через API (см. эндпойнты ниже) и смотрите метрики/дашборды.

## Локальная разработка (Go)

Требуется Go 1.24+ и Postgres.

1. Поднимите Postgres (`docker-compose up db` или свой инстанс). Значения по умолчанию в dev-конфиге: host `localhost`,
   port `6132`, user/pass `postgres`, db `uptime_monitor_db`.
2. Убедитесь, что есть `config/config.prod.yaml` (можно скопировать из `config/config.dev.yaml` и поправить).
3. Запустите сервис:
   ```sh
   go run ./cmd
   ```
   По умолчанию слушает `0.0.0.0:8080`.

## API (HTTP, JSON)

Базовый путь: `/api`

| Метод | Путь | Тело/параметры | Описание | Ответ |
| --- | --- | --- | --- | --- |
| POST | `/api/url` | JSON `{"url": "<https://...>"}` | Добавить URL | 201 Created + данные URL |
| DELETE | `/api/url` | JSON `{"url": "<https://...>"}` | Удалить URL | 204 No Content |
| POST | `/api/url/activate` | JSON `{"url": "<https://...>"}` | Активировать URL | 204 No Content |
| POST | `/api/url/deactivate` | JSON `{"url": "<https://...>"}` | Деактивировать URL | 204 No Content |
| GET | `/api/url` | — | Список URL | 200 + массив URL |
| GET | `/api/url/:id/history` | path `id` | История (статус, задержка, время) | 200 + массив history |
| GET | `/api/metrics` | — | Метрики Prometheus | 200 text/plain |

OpenAPI-спека: `openapi.yml` в корне.

## Конфигурация

Файл: `config/config.prod.yaml` (Viper). Основные поля:

- `Server.host`, `Server.port` — адрес/порт HTTP.
- `Postgres.*` — параметры подключения (`sslMode` по умолчанию `disable`).
- `Monitor.intervalSeconds`, `requestTimeoutSeconds`, `historyLimit`, `maxConcurrency`.

## Наблюдаемость

- Метрики (Prometheus):
    - `monitor_urls_in_flight` — текущие параллельные проверки.
    - `monitor_last_status_code{url_id,url}` — последний HTTP-статус по URL.
    - `monitor_last_latency_ms{url_id,url}` — последняя задержка по URL.
    - `monitor_url_active{url_id,url}` — 1 если URL активен, 0 если выключен.
    - HTTP-миддлвар: `http_requests_total`, `http_request_duration_seconds`, `http_requests_in_flight`.
- Логи: stdout (zap). Promtail/Loki подключены через compose (`promtail-config.yaml`, `loki-config.yaml`).
- Grafana: провиженинг датасорса в `grafana/datasources.yaml` (admin/admin).

## Хранение и миграции

- Миграции в `db/migrations`, применяются на старте через Goose.
- История запросов хранится в таблице history (см. `internal/entities/history.go`).

## Структура проекта

- `cmd/` — входная точка.
- `config/` — модели и примеры конфигов.
- `internal/monitor/` — планировщик и логика опросов.
- `internal/repository/postgres/` — работа с Postgres + миграции.
- `internal/delivery/http/` — HTTP-сервер, роуты, хендлеры, middleware.
- `internal/metrics/` — Prometheus-метрики.
- `openapi.yml` — контракт API.

## Полезные команды

- Формат/линт: `make fmt` / `make lint`.
- Запуск стека: `docker-compose up --build`.
- Быстрый просмотр метрик: `curl http://localhost:8080/api/metrics | head`.
