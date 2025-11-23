# Makefile for the project
# Provides common convenience targets and default environment variables.
.DEFAULT_GOAL := help

# Directories / binary
APP_DIR := .
BINARY ?= app

# Environment defaults (can be overridden on the command line)
REDIS_ADDR ?= localhost:6379
GATEWAY_ADDR ?= http://localhost:8888

# Go tool
GO ?= go
BUILD_FLAGS ?=
TEST_FLAGS ?=

# Metadata
VERSION ?= dev

.PHONY: help run build test deps clean version

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Common targets:"
	@echo "  run       - run the application with REDIS_ADDR and GATEWAY_ADDR set (default values provided)"
	@echo "  build     - build the binary ($(BINARY))"
	@echo "  test      - run unit tests"
	@echo "  deps      - tidy modules (go mod tidy)"
	@echo "  clean     - remove built binary"
	@echo "  version   - print build version"

# Run the app with configured environment variables.
# You can override e.g. `make run REDIS_ADDR=redis:6379 GATEWAY_ADDR=http://gateway:8888`
run:
	@docker compose up -d
	REDIS_ADDR=$(REDIS_ADDR) GATEWAY_ADDR=$(GATEWAY_ADDR) $(GO) run $(APP_DIR)

# Build the project binary
build:
	$(GO) build -o $(BINARY) $(BUILD_FLAGS) $(APP_DIR)

# Run tests
test:
	$(GO) test ./... $(TEST_FLAGS)

# Ensure module dependencies are tidy
deps:
	$(GO) mod tidy

# Clean built artifacts
clean:
	@docker compose down
	@rm -f $(BINARY)

# Print version
version:
	@echo "version: $(VERSION)"
