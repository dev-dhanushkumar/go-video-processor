.PHONY: help build run test clean docker-build docker-run migrate-up migrate-down install-deps

# Variables
APP_NAME=video-processor
MAIN_PATH=cmd/server/main.go
BUILD_DIR=bin
DOCKER_IMAGE=video-processor:latest

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install-deps: ## Install Go dependencies
	go mod download
	go mod tidy

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run: ## Run the application
	go run $(MAIN_PATH)

dev: ## Run in development mode with auto-reload (requires air)
	air

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@for file in migrations/*_*.up.sql; do \
		echo "Applying $$file"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f $$file; \
	done

migrate-down: ## Run database migrations down
	@echo "Reverting migrations..."
	@for file in $$(ls -r migrations/*_*.down.sql); do \
		echo "Reverting $$file"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f $$file; \
	done

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

lint: ## Run linter
	golangci-lint run

format: ## Format code
	go fmt ./...

setup: install-deps ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env
	@mkdir -p storage/uploads storage/processed storage/thumbnails
	@echo "Setup complete! Please configure .env file"

.DEFAULT_GOAL := help
