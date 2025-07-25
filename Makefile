# i3 Config Generator - Makefile

BINARY_NAME := i3-config-generator
BUILD_DIR := build
INSTALL_DIR := $(HOME)/.local/bin
CONFIG_DIR := $(HOME)/.config/i3-config-generator

.PHONY: all build test test-binary clean install setup-config check fmt vet deps help

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "✓ All tests passed"

# Test the built binary
test-binary: build
	@echo "Testing binary with different layouts..."
	@mkdir -p /tmp/i3-config-test
	$(BUILD_DIR)/$(BINARY_NAME) --version
	$(BUILD_DIR)/$(BINARY_NAME) --layout two_mon --output /tmp/i3-config-test/two_mon
	$(BUILD_DIR)/$(BINARY_NAME) --layout one_mon --output /tmp/i3-config-test/one_mon
	$(BUILD_DIR)/$(BINARY_NAME) --layout no_mon --output /tmp/i3-config-test/no_mon
	@echo "✓ Binary tests completed"
	@echo "Generated configs:"
	@ls -l /tmp/i3-config-test/
	@rm -rf /tmp/i3-config-test

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ Vet checks passed"

# Run all quality checks
check: fmt vet test
	@echo "✓ All quality checks passed"

# Install binary
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "✓ Installed: $(INSTALL_DIR)/$(BINARY_NAME)"

# Set up configuration
setup-config:
	@echo "Setting up configuration..."
	@mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_DIR)/config.yaml ]; then \
		cp config.yaml $(CONFIG_DIR)/config.yaml; \
		echo "✓ Config installed: $(CONFIG_DIR)/config.yaml"; \
	else \
		echo "ℹ Config already exists: $(CONFIG_DIR)/config.yaml"; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	@echo "✓ Dependencies downloaded"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

# Show help
help:
	@echo "i3 Config Generator - Available targets:"
	@echo ""
	@echo "  build        - Build the binary"
	@echo "  test         - Run all tests"
	@echo "  test-binary  - Test the built binary with different layouts"
	@echo "  check        - Run all quality checks (fmt + vet + test)"
	@echo "  install      - Install binary to ~/.local/bin"
	@echo "  setup-config - Set up configuration directory"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download dependencies"
	@echo "  help         - Show this help message"