.PHONY: build clean test install help

# Variables
BINARY_NAME=pdf2img
BINARY_PATH=./cmd/pdf2img
GO=go
LDFLAGS=-ldflags "-s -w"

help:
	@echo "pdf2img - PDF to Image Converter"
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build the CLI application"
	@echo "  build-dev   - Build with debug symbols"
	@echo "  clean       - Remove compiled binaries"
	@echo "  install     - Install globally"
	@echo "  test        - Run tests"
	@echo "  deps        - Download dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"

deps:
	$(GO) mod download
	$(GO) mod verify

build: deps
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "✓ Built $(BINARY_NAME)"

build-dev: deps
	$(GO) build -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "✓ Built $(BINARY_NAME) (with debug symbols)"

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	@echo "✓ Cleaned"

install: build
	$(GO) install $(BINARY_PATH)
	@echo "✓ Installed globally"

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...
	@echo "✓ Code formatted"

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
	golangci-lint run ./...

all: clean build
	@echo "✓ All done"
