# Makefile for clitools

# Variables
BINARY_NAME=clitools
BUILD_DIR=build
MAIN_FILE=main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Built $(BINARY_NAME) successfully!"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	@echo "Built all platform binaries in $(BUILD_DIR)/"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(GOPATH)/bin/"

# Development: build and run with example
.PHONY: dev
dev: build
	@echo "Running development version..."
	@./$(BINARY_NAME) --help

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  deps       - Install dependencies"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install binary to GOPATH/bin"
	@echo "  dev        - Build and show help"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  help       - Show this help"
