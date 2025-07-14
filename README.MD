# USDT Rate Service (GRPC)

gRPC-сервис на Go, который получает курс USDT с биржи Grinex и сохраняет его в PostgreSQL.

## 🚀 Возможности

- Получение курса USDT (ask, bid)
- Хранение истории курсов в PostgreSQL
- Healthcheck
- Graceful shutdown
- Unit-тесты
- Конфигурация через `.env` и флаги
- Логирование через zap
- Поддержка миграций (Goose)
- Docker/Docker Compose
- Makefile-команды

---

## 🧰 Требования

- Go >= 1.22
- Docker + Docker Compose
- Make

---

## ⚙️ Конфигурация

Создай файл `.env` в корне проекта:

```env
GRPC_PORT=50051
GRINEX_URL=https://grinex.io/api/v2/depth
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=usdt_db
DB_HOST=postgres
DB_PORT=5432
```

## 🧪 Команды Makefile

| Команда             | Описание                          |
|---------------------|-----------------------------------|
| `make build`        | Сборка приложения                 |
| `make run`          | Запуск локального приложения      |
| `make docker-build` | Сборка Docker-образа              |
| `make migrate-up`   | Применение миграций (goose)       |
| `make migrate-down` | Откат миграций                    |
| `make lint`         | Проверка линтером (golangci-lint) |
| `make test`         | Прогон тестов                     |


---

## 🐳 Быстрый запуск (через Docker)

```bash
git clone https://github.com/iskanarynbaev/usdt-rate-service.git
cd usdt-rate-service

# Сборка и запуск всего
make build
docker-compose up -d
make migrate-up

# Запуск приложения внутри контейнера
docker-compose run --rm app ./app
