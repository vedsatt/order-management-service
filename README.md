# Order Management Service

gRPC, gRPC-Gateway микросервис для управления заказами с PostgreSQL хранилищем и кэшем на основе Redis.

## Установка и запуск

### Требования
- Docker
- docker-compose

### Быстрый старт

```bash
# Клонирование репозитория
git clone https://gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478.git
cd order-management-service

# Установка зависимостей
make deps

# Запуск сервера
make docker-build
```

## Конфигурация

Сервис настраивается через файл .env (подробнее в ```./config/env.example```):

| Переменная         | По умолчанию | Описание                 |
|--------------------|--------------|--------------------------|
| `GRPC_PORT`        | `50051`      | Порт gRPC сервера        |
| `GATEWAY_PORT`     | `8080`       | Порт gRPC Gateway        |
| `ENV`              | `prod`       | Окружение (`dev`/`prod`) |
| Конфигурация базы данных:                                    |
| `POSTGRES_HOST`    | `postgres`   |                          |
| `POSTGRES_VERSION` | `15-alpine`  |                          |
| `POSTGRES_DB`      | `postgres`   |                          |
| `POSTGRES_USER`    | `postgres`   |                          |
| `POSTGRES_PASSWORD`| `postgres`   |                          |
| `POSTGRES_PORT`    | `5432`       |                          |
| Конфигурация кэша:                                           |
| `REDIS_HOST`       | `redis`      |                          |
| `REDIS_VERSION`    | `8.0-alpine` |                          |
| `REDIS_PASSWORD`   | `redis`      |                          |
| `REDIS_PORT`       | `6379`       |                          |
| `REDIS_MAX_MEMORY` | `256mb`      |                          |
|--------------------|--------------|--------------------------|

## Makefile
Список и описание функционала всех доступных команд:
```bash
make help
```
