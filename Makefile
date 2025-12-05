.PHONY: help build run test migrate-up migrate-down docker-up docker-down clean swagger

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@go build -o bin/urlshortener cmd/api/main.go
	@echo "Build complete: bin/urlshortener"

run: ## Run the application
	@echo "Starting URL Shortener Service..."
	@go run cmd/api/main.go

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@$(HOME)/go/bin/swag init -g cmd/api/main.go
	@echo "Swagger docs generated"

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@docker exec -it url_shortener_db psql -U postgres -d url_shortener -f /migrations/000001_create_urls_table.up.sql
	@echo "Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@docker exec -it url_shortener_db psql -U postgres -d url_shortener -f /migrations/000001_create_urls_table.down.sql
	@echo "Rollback complete"

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Docker containers started"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Docker containers stopped"

docker-logs: ## Show Docker container logs
	@docker-compose logs -f

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

dev: docker-up swagger run ## Start development environment (Docker + Swagger + Run)

.DEFAULT_GOAL := help
