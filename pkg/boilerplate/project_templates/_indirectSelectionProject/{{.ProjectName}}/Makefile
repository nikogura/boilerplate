.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: deps
deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest

.PHONY: proto
proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/example-indirect-selection/service.proto

.PHONY: generate
generate: proto ## Run go generate and protobuf generation
	@echo "Running go generate..."
	go generate ./...

.PHONY: tidy
tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	go mod tidy

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: test
test: ## Run tests with race detection and coverage
	@echo "Running tests..."
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

.PHONY: coverage
coverage: test ## Show test coverage
	@echo "Showing coverage..."
	go tool cover -html=coverage.out

.PHONY: ci
ci: deps tidy generate lint test ## Run full CI pipeline

.PHONY: build
build: generate ## Build the binary
	@echo "Building binary..."
	go build -o bin/example-indirect-selection .

.PHONY: run
run: build ## Run the service
	@echo "Starting service..."
	./bin/example-indirect-selection server

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out