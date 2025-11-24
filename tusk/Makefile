.PHONY: build test test-e2e clean release install help

# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X github.com/sboy99/projects.go/tusk/internal/version.BuildVersion=$(VERSION)
LDFLAGS += -X github.com/sboy99/projects.go/tusk/internal/version.BuildCommit=$(COMMIT)
LDFLAGS += -X github.com/sboy99/projects.go/tusk/internal/version.BuildTime=$(BUILD_TIME)

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build the tusk binary"
	@echo "  test       - Run unit tests"
	@echo "  test-e2e   - Run end-to-end tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  release    - Build release binaries for multiple platforms"
	@echo "  install    - Install tusk to GOPATH/bin"
	@echo "  deps       - Download dependencies"

# Build the binary
build:
	@mkdir -p bin
	@./scripts/build.sh

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run e2e tests
test-e2e:
	@echo "Running e2e tests..."
	@go test -v ./test/...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ dist/
	@go clean

# Build release binaries
release:
	@echo "Building release binaries..."
	@./scripts/release.sh

# Install to GOPATH/bin
install:
	@echo "Installing tusk..."
	@go install -ldflags "$(LDFLAGS)" ./cmd/tusk
	@echo "Installation complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

