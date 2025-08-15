# Pit - A tiny, educational Git implementation in Go

# Build configuration
BINARY_NAME=pit
BUILD_DIR=bin
SOURCE_DIR=cmd/pit
MAIN_FILE=$(SOURCE_DIR)/main.go

# Go configuration
GO_VERSION=1.21
GOFLAGS=-ldflags="-s -w"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(MAIN_FILE) $(shell find . -name "*.go" -not -path "./vendor/*")
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Built $(BINARY_NAME) successfully"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@echo "Clean completed"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Development build (with debug symbols)
.PHONY: dev
dev:
	@echo "Building development version..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Run the binary
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(MAIN_FILE)

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Quality checks (format, lint, test)
.PHONY: quality
quality: fmt test
	@echo "All quality checks passed!"

# Git integration test
.PHONY: git-test
git-test: build
	@echo "Running Git compatibility tests..."
	@mkdir -p /tmp/pit-git-test
	@cd /tmp/pit-git-test && ../../../$(shell pwd)/$(BUILD_DIR)/$(BINARY_NAME) init
	@echo "Git compatibility test completed"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Build the project (default)"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Install dependencies"
	@echo "  dev           - Development build with debug symbols"
	@echo "  run           - Build and run the binary"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  quality       - Run all quality checks (fmt + test)"
	@echo "  git-test      - Run Git compatibility tests"
	@echo "  help          - Show this help"

# Prevent make from treating files with these names as targets
.PHONY: $(BUILD_DIR)/$(BINARY_NAME)