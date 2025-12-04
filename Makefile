.PHONY: help build run test clean docker-build docker-run lint fmt vet

# Variables
SERVICE_NAME=redirect-service
DOCKER_IMAGE=gcr.io/PROJECT_ID/$(SERVICE_NAME)
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the service binary
	@echo "Building $(SERVICE_NAME)..."
	@go build -o bin/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

run: ## Run the service locally
	@echo "Running $(SERVICE_NAME)..."
	@go run ./cmd/$(SERVICE_NAME)

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ dist/ coverage.* *.log

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):latest .

docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	@docker run -p 8081:8081 --env-file .env $(DOCKER_IMAGE):latest

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

deps-upgrade: ## Upgrade dependencies
	@echo "Upgrading dependencies..."
	@go get -u ./...
	@go mod tidy

.DEFAULT_GOAL := help
