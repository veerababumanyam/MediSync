# MediSync Makefile
# Build, test, and run commands for the MediSync AI Agent Core

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Binary names
API_BINARY=medisync-api
MIGRATE_BINARY=medisync-migrate

# Main packages
API_MAIN=./cmd/api
MIGRATE_MAIN=./cmd/migrate

# Build directory
BUILD_DIR=./bin

# Docker
DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE_FILE=docker-compose.yml

# Linter
LINTER=golangci-lint

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test lint run migrate docker-up docker-down fmt deps help

all: deps build

## build: Build all binaries
build:
	@echo "$(GREEN)Building binaries...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(API_BINARY) $(API_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(MIGRATE_BINARY) $(MIGRATE_MAIN)
	@echo "$(GREEN)Build complete!$(NC)"

## build-api: Build only the API binary
build-api:
	@echo "$(GREEN)Building API binary...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(API_BINARY) $(API_MAIN)
	@echo "$(GREEN)API build complete!$(NC)"

## build-migrate: Build only the migrate binary
build-migrate:
	@echo "$(GREEN)Building migrate binary...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(MIGRATE_BINARY) $(MIGRATE_MAIN)
	@echo "$(GREEN)Migrate build complete!$(NC)"

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Clean complete!$(NC)"

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests complete!$(NC)"

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "$(GREEN)Generating coverage report...$(NC)"
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## test-unit: Run only unit tests
test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GOTEST) -v -race -short ./tests/...
	@echo "$(GREEN)Unit tests complete!$(NC)"

## test-integration: Run only integration tests
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GOTEST) -v -race ./tests/integration/...
	@echo "$(GREEN)Integration tests complete!$(NC)"

## test-contract: Run only contract tests
test-contract:
	@echo "$(GREEN)Running contract tests...$(NC)"
	$(GOTEST) -v -race ./tests/contract/...
	@echo "$(GREEN)Contract tests complete!$(NC)"

## lint: Run linters
lint:
	@echo "$(GREEN)Running linters...$(NC)"
	$(LINTER) run ./...
	@echo "$(GREEN)Linting complete!$(NC)"

## lint-fix: Run linters with auto-fix
lint-fix:
	@echo "$(GREEN)Running linters with auto-fix...$(NC)"
	$(LINTER) run --fix ./...
	@echo "$(GREEN)Linting with auto-fix complete!$(NC)"

## fmt: Format Go code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@echo "$(GREEN)Formatting complete!$(NC)"

## run: Run the API server
run:
	@echo "$(GREEN)Starting API server...$(NC)"
	$(GOCMD) run $(API_MAIN)

## migrate: Run database migrations
migrate:
	@echo "$(GREEN)Running database migrations...$(NC)"
	$(GOCMD) run $(MIGRATE_MAIN)

## migrate-up: Run migrations up
migrate-up:
	@echo "$(GREEN)Running migrations up...$(NC)"
	$(GOCMD) run $(MIGRATE_MAIN) up

## migrate-down: Rollback migrations
migrate-down:
	@echo "$(YELLOW)Rolling back migrations...$(NC)"
	$(GOCMD) run $(MIGRATE_MAIN) down

## migrate-status: Check migration status
migrate-status:
	@echo "$(GREEN)Checking migration status...$(NC)"
	$(GOCMD) run $(MIGRATE_MAIN) status

## docker-up: Start Docker containers
docker-up:
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "$(GREEN)Docker containers started!$(NC)"

## docker-down: Stop Docker containers
docker-down:
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down
	@echo "$(GREEN)Docker containers stopped!$(NC)"

## docker-logs: View Docker container logs
docker-logs:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs -f

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies downloaded!$(NC)"

## deps-update: Update all dependencies
deps-update:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	$(GOMOD) tidy
	$(GOGET) -u ./...
	@echo "$(GREEN)Dependencies updated!$(NC)"

## verify: Verify the codebase is ready for commit
verify: fmt lint test
	@echo "$(GREEN)Verification complete!$(NC)"

## install-linter: Install golangci-lint
install-linter:
	@echo "$(GREEN)Installing golangci-lint...$(NC)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
	@echo "$(GREEN)golangci-lint installed!$(NC)"

## env-example: Create example .env file
env-example:
	@echo "$(GREEN)Creating example .env file...$(NC)"
	@cp .env.example .env 2>/dev/null || echo "DATABASE_URL=postgres://medisync:medisync@localhost:5432/medisync?sslmode=disable\nREDIS_URL=redis://localhost:6379\nKEYCLOAK_URL=http://localhost:8081\nOPA_URL=http://localhost:8181\nNATS_URL=nats://localhost:4222\nLLM_PROVIDER=ollama\nLLM_MODEL=llama4\nLLM_API_KEY=\nSERVER_PORT=:8080\nENVIRONMENT=development\nLOG_LEVEL=debug" > .env
	@echo "$(GREEN).env file created!$(NC)"

## help: Show this help message
help:
	@echo "MediSync Makefile Commands:"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
	@echo ""
	@echo "Examples:"
	@echo "  make build          # Build all binaries"
	@echo "  make test           # Run all tests"
	@echo "  make lint           # Run linters"
	@echo "  make run            # Start the API server"
	@echo "  make docker-up      # Start Docker containers"
	@echo "  make verify         # Run fmt, lint, and test"

# Default target
.DEFAULT_GOAL := help
