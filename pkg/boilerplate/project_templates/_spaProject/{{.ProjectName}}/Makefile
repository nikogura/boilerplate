IMAGE_NAME = "{{.ProjectName}}"

.PHONY: build clean test lint generate deps tidy ci run help

# Default target
all: build

# Generate build info for UI
build-info:
	@echo "Generating build info..."
	@BUILD_VERSION="${SEMVER:-local}" && \
	mkdir -p pkg/ui/static && \
	echo "{\"version\":\"$$BUILD_VERSION\"}" > pkg/ui/static/build-info.json

# Build {{.ProjectName}} binary
build: build-info
	@echo "Building {{.ProjectName}} binary..."
	go build -o {{.ProjectName}}

# Clean built files
clean:
	@echo "Cleaning binaries..."
	rm -f {{.ProjectName}} coverage.out

# Install dependencies and tools
deps:
	@echo "Installing tools and dependencies..."
	go mod download
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	fi

# Generate mocks and assets
generate:
	@echo "Running go generate..."
	go generate ./...

# Run linters
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Run tests with coverage
test:
	@echo "Running tests with race detection and coverage..."
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

# Tidy dependencies
tidy:
	@echo "Running go mod tidy..."
	go mod tidy

# Full CI pipeline
ci: deps tidy generate lint test
	@echo "✅ CI pipeline completed successfully"

# Run the server for development
run: build
	@echo "Starting {{.ProjectName}} server..."
	./{{.ProjectName}} server

# Run the server with debug logging
run-debug: build
	@echo "Starting {{.ProjectName}} server with debug logging..."
	{{.ProjectEnvPrefix}}_LOG_LEVEL=debug ./{{.ProjectName}} server

# Run with sample OIDC configuration (requires setting up OAuth2 credentials)
run-with-auth: build
	@echo "Starting {{.ProjectName}} server with OIDC authentication..."
	@echo "Note: Set {{.ProjectEnvPrefix}}_OIDC_CLIENT_ID and {{.ProjectEnvPrefix}}_OIDC_CLIENT_SECRET"
	{{.ProjectEnvPrefix}}_LOG_LEVEL=debug \
	{{.ProjectEnvPrefix}}_OIDC_CLIENT_ID="${OIDC_CLIENT_ID}" \
	{{.ProjectEnvPrefix}}_OIDC_CLIENT_SECRET="${OIDC_CLIENT_SECRET}" \
	{{.ProjectEnvPrefix}}_OIDC_REDIRECT_URL="http://localhost:9999/auth/callback" \
	./{{.ProjectName}} server

# Show coverage report
coverage: test
	@echo "Opening coverage report..."
	go tool cover -html=coverage.out

# Docker build
docker-build: build
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build {{.ProjectName}} binary (default)"
	@echo "  clean        - Clean built files"
	@echo "  deps         - Install dependencies and tools"
	@echo "  generate     - Run go generate"
	@echo "  lint         - Run golangci-lint"
	@echo "  test         - Run tests with race detection and coverage"
	@echo "  tidy         - Run go mod tidy"
	@echo "  ci           - Run full CI pipeline (deps + tidy + generate + lint + test)"
	@echo "  run          - Build and run the server"
	@echo "  run-debug    - Build and run with debug logging"
	@echo "  run-with-auth - Build and run with OIDC auth (requires env vars)"
	@echo "  coverage     - Show test coverage report"
	@echo "  docker-build - Build Docker image"
	@echo "  help         - Show this help message"