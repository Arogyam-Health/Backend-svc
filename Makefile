.PHONY: help build run test test-integration test-unit clean pre-commit fmt vet lint docker-build docker-run

# Variables
APP_NAME=backend-service
DOCKER_IMAGE=instagram-backend
DOCKER_TAG=latest
GO=go
GOFLAGS=-v

# Default target
help:
	@echo "Available targets:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make test             - Run all tests"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make test-unit        - Run unit tests only"
	@echo "  make pre-commit       - Run pre-commit checks (required before commit)"
	@echo "  make fmt              - Format code"
	@echo "  make vet              - Run go vet"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	$(GO) build -o app ./cmd/main.go
	@echo "Build complete: ./app"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GO) run ./cmd/main.go

# Run all tests
test:
	@echo "Running all tests..."
	$(GO) test $(GOFLAGS) -cover ./...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	$(GO) test $(GOFLAGS) ./tests/integration/...

# Run unit tests only (excluding integration tests)
test-unit:
	@echo "Running unit tests..."
	$(GO) test $(GOFLAGS) ./internal/...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Lint code (requires golangci-lint to be installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping..."; \
		echo "Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

# Pre-commit checks - MUST pass before committing
pre-commit: fmt vet test build
	@echo ""
	@echo "✅ Pre-commit checks passed!"
	@echo ""
	@echo "📝 Reminder:"
	@echo "  1. Update CHANGELOG.md with your changes"
	@echo "  2. Ensure your commit message follows the format:"
	@echo "     <type>: <description>"
	@echo ""
	@echo "  Types: feat, fix, test, docs, refactor, chore"
	@echo ""
	@echo "Example: feat: add media deletion endpoint"
	@echo ""

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f app
	rm -f test_token.json
	rm -f token.json
	$(GO) clean
	@echo "Clean complete"

# Docker build
docker-build:
	@echo "Building Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built successfully"

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -d \
		--name $(APP_NAME) \
		--env-file .env \
		-p 8080:8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Container started. Check with: docker ps"

# Docker stop and remove
docker-stop:
	@echo "Stopping and removing container..."
	-docker stop $(APP_NAME)
	-docker rm $(APP_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies installed"

# Check if code is formatted
check-fmt:
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "❌ Code is not formatted. Run 'make fmt'"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "✅ Code is properly formatted"; \
	fi

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Quick check (faster than pre-commit, for rapid iteration)
quick-check: fmt vet
	@echo "✅ Quick checks passed!"
