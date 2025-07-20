.PHONY: help build run test clean deps lint

# Переменные
BINARY_NAME=go_news_server
BUILD_DIR=build
MAIN_PATH=cmd/go_news_server/main.go

# Цвета для вывода
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

deps: ## Установить зависимости
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	go mod download
	go mod tidy

build: ## Собрать приложение
	@echo "$(GREEN)Сборка приложения...$(NC)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: ## Запустить приложение
	@echo "$(GREEN)Запуск приложения...$(NC)"
	go run $(MAIN_PATH)

dev: ## Запустить в режиме разработки
	@echo "$(GREEN)Запуск в режиме разработки...$(NC)"
	STAGE_STATUS=dev go run $(MAIN_PATH)

test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	go test ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(GREEN)Запуск тестов с покрытием...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint: ## Проверить код линтером
	@echo "$(GREEN)Проверка кода...$(NC)"
	golangci-lint run

clean: ## Очистить сборки
	@echo "$(GREEN)Очистка...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

docker-build: ## Собрать Docker образ
	@echo "$(GREEN)Сборка Docker образа...$(NC)"
	docker build -t $(BINARY_NAME) .

docker-run: ## Запустить в Docker
	@echo "$(GREEN)Запуск в Docker...$(NC)"
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME)

setup-db: ## Настройка базы данных
	@echo "$(GREEN)Настройка базы данных...$(NC)"
	@echo "Создайте базу данных и примените схему:"
	@echo "createdb news_db"
	@echo "psql -d news_db -f database/schema.sql"

setup-env: ## Настройка переменных окружения
	@echo "$(GREEN)Настройка переменных окружения...$(NC)"
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Файл .env создан из env.example"; \
		echo "Отредактируйте .env файл под ваши настройки"; \
	else \
		echo "Файл .env уже существует"; \
	fi

install: deps setup-env ## Полная установка
	@echo "$(GREEN)Установка завершена!$(NC)"
	@echo "Не забудьте настроить базу данных: make setup-db"