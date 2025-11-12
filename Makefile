BINARY_NAME=order-service
GRPC_PORT=50051

.PHONY: all build up down exec logs restart lint deps grpc grpc-gateway help

all: docker-build

DCE = docker-compose --env-file config/.env

build:
	$(DCE) up -d --build

up:
	$(DCE) up -d

down:
	$(DCE) down

exec:
	docker exec -it order_management_service sh

logs:
	$(DCE) logs -f app

restart:
	$(DCE) restart app

lint:
	@echo "Running linters..."
	golangci-lint run

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

grpc:
	@echo "Generating gRPC..."
	protoc -I . --go_out . --go-grpc_out . api/order.proto

gateway:
	@echo "Generating gRPC-Gateway..."
	protoc -I . --grpc-gateway_out . api/order.proto 

test:
	go test ./... -v

help:
	@echo "Available commands:"
	@echo "  make build         - Собрать проект"
	@echo "  make up            - Запустить сервер"
	@echo "  make down          - Остановить сервер"
	@echo "  make exec          - Подключиться к shell контейнера"
	@echo "  make restart       - Перезапустить сервер"
	@echo "  make logs          - Вывести логи"
	@echo "  make grpc          - Сгенерировать gRPC"
	@echo "  make gateway       - Сгенерировать gRPC-Gateway"
	@echo "  make lint          - Проверить код линтером"
	@echo "  make deps          - Установить зависимости"
	@echo "  make test          - Запустить тесты"
	@echo "  make help          - Показать эту справку"