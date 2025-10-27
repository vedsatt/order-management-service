BINARY_NAME=order-service
GRPC_PORT=50051

.PHONY: all build run clean lint help

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) ./cmd/server

build_and_run:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) ./cmd/server
	@echo "Running ./bin/$(BINARY_NAME)..."
		bin/$(BINARY_NAME)

run:
	@echo "Starting server on port $(GRPC_PORT)..."
	go run ./cmd/server

clean:
	@echo "Cleaning up..."
	rm -rf bin/

lint:
	@echo "Running linters..."
	golangci-lint run

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

help:
	@echo "Available commands:"
	@echo "  make build         - Собрать проект"
	@echo "  make build and run - Собрать и запустить сервер"
	@echo "  make run           - Запустить сервер"
	@echo "  make lint          - Проверить код линтером"
	@echo "  make deps          - Установить зависимости"
	@echo "  make clean         - Очистить билды"
	@echo "  make help          - Показать эту справку"