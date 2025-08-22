.PHONY: deps generate lint test ci run tidy clean build

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest

# Generate mocks and other generated code
generate:
	@echo "Generating code..."
	go generate ./...

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Run tests with race detection and coverage
test:
	@echo "Running tests..."
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

# Run full CI pipeline
ci: tidy generate lint test
	@echo "CI pipeline completed successfully"

# Build the application
build:
	@echo "Building application..."
	mkdir -p bin
	go build -o bin/{{.ProjectName}} .

# Run the service
run: build
	@echo "Starting {{.ProjectName}} service..."
	@echo "Metrics server will be available at http://localhost:{{.DefaultServerPort}}/metrics"
	@echo "Health endpoint: http://localhost:{{.DefaultServerPort}}/healthz"
	@echo "Ready endpoint: http://localhost:{{.DefaultServerPort}}/readyz"
	./bin/{{.ProjectName}}

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin coverage.out