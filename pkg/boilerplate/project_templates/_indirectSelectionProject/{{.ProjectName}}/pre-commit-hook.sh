#!/usr/bin/env bash

set -e

# Run go mod tidy
echo "Running go mod tidy..."
go mod tidy

# Run go vet
echo "Running go vet..."
go vet ./...

# Run golint if available
if command -v golint &> /dev/null; then
    echo "Running golint..."
    golint ./...
fi

# Run go fmt check
echo "Checking go fmt..."
GOFMT_FILES=$(gofmt -l .)
if [[ -n ${GOFMT_FILES} ]]; then
    echo "gofmt check failed for:"
    echo "${GOFMT_FILES}"
    exit 1
fi

echo "All checks passed!"