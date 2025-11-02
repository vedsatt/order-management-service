# Order Management Service

gRPC микросервис для управления заказами с in-memory хранилищем.

## Установка и запуск

### Требования
- Go 1.21 или выше
- Протокол компилятор (protoc)

### Быстрый старт

```bash
# Клонирование репозитория
git clone https://gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478.git
cd order-management-service

# Установка зависимостей
make deps

# Запуск сервера
make run
```

## Конфигурация

Сервис настраивается через файл .env (подробнее в ```./config/env.example```):

| Переменная    | По умолчанию | Описание                 |
|---------------|--------------|--------------------------|
| `GRPC_PORT`   | `50051`      | Порт gRPC сервера        |
| `GATEWAY_PORT`| `prod`       | Порт gRPC Gateway        |
| `ENV`         | `prod`       | Окружение (`dev`/`prod`) |


## Makefile
Список и описание функционала всех доступных команд:
```bash
make help
```
