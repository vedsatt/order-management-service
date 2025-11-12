# Order Management Service

Микросервис для управления заказами с gRPC API и HTTP Gateway.

## Требования
- **Docker** 20.10+
- **docker-compose** 2.0+

## Технологии

### Backend
- **Go 1.24** - основной язык программирования
- **gRPC** - высокопроизводительный RPC фреймворк
- **gRPC Gateway** - HTTP/JSON прокси для gRPC API
- **Zap** - структурированное логирование

### База данных
- **PostgreSQL** - основная реляционная БД
- **Pgx** - драйвер PostgreSQL для Go
- **Squirrel** - SQL builder для построения запросов
- **Golang Migrate** - миграции базы данных

### Кеширование
- **Redis** - in-memory кеш
- **go-redis** - клиент Redis для Go

### Инфраструктура
- **Docker** - контейнеризация
- **Docker Compose** - оркестрация сервисов
- **Make** - автоматизация задач

### Разработка
- **Testify** - фреймворк для тестирования
- **Cleanenv** - управление конфигурацией
- **Protocol Buffers** - сериализация данных и определение API

## Установка и запуск

```bash
# Клонирование репозитория
git clone https://gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478.git
cd order-management-service

# Запуск сервера
make build
```

## Структура проекта
```
.
├── Dockerfile
├── Makefile
├── README.md
├── api
│   └── order.proto
├── cmd
│   ├── migrate
│   │   └── main.go
│   └── server
│       └── main.go
├── config
│   └── env.example
├── docker-compose.yml
├── go.mod
├── go.sum
├── google
│   └── api
│       ├── annotations.proto
│       ├── field_behavior.proto
│       ├── http.proto
│       └── httpbody.proto
├── internal
│   ├── config
│   │   └── config.go
│   ├── repository
│   │   ├── cache
│   │   │   └── order_cache.go
│   │   ├── database
│   │   │   └── order_repo.go
│   │   └── order_repository.go
│   ├── service
│   │   ├── order.go
│   │   └── order_test.go
│   └── transport
│       ├── gateway.go
│       ├── grpc_order_server.go
│       └── grpc_order_server_test.go
├── migrations
│   ├── 001_create_order_table.down.sql
│   └── 001_create_order_table.up.sql
└── pkg
    ├── api
    │   └── test
    │       ├── order.pb.go
    │       ├── order.pb.gw.go
    │       └── order_grpc.pb.go
    └── logger
        └── logger.go
```

## Конфигурация

Сервис настраивается через файл .env (подробнее в ```./config/env.example```):

| Переменная         | По умолчанию | Описание                 |
|--------------------|--------------|--------------------------|
| `GRPC_PORT`        | `50051`      | Порт gRPC сервера        |
| `GATEWAY_PORT`     | `8080`       | Порт gRPC Gateway        |
| `ENV`              | `prod`       | Окружение (`dev`/`prod`) |
| **Конфигурация базы данных:**                                |
| `POSTGRES_HOST`    | `postgres`   |                          |
| `POSTGRES_VERSION` | `15-alpine`  |                          |
| `POSTGRES_DB`      | `postgres`   |                          |
| `POSTGRES_USER`    | `postgres`   |                          |
| `POSTGRES_PASSWORD`| `postgres`   |                          |
| `POSTGRES_PORT`    | `5432`       |                          |
| **Конфигурация кэша:**                                       |
| `REDIS_HOST`       | `redis`      |                          |
| `REDIS_VERSION`    | `8.0-alpine` |                          |
| `REDIS_PASSWORD`   | `redis`      |                          |
| `REDIS_PORT`       | `6379`       |                          |
| `REDIS_MAX_MEMORY` | `256mb`      |                          |

## Makefile
Список и описание функционала всех доступных команд:
```bash
make help
```
