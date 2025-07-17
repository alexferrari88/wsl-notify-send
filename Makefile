# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=wsl-notify-send
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build targets
.PHONY: all build build-windows build-linux clean test coverage deps help

all: clean deps test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v .

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v .

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)
	rm -f coverage.out
	rm -f coverage.html

test:
	$(GOTEST) -v ./...

test-race:
	$(GOTEST) -v -race ./...

coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out

deps:
	$(GOMOD) download
	$(GOMOD) tidy

fmt:
	$(GOCMD) fmt ./...

lint:
	golangci-lint run

vet:
	$(GOCMD) vet ./...

# Install development dependencies
install-dev:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all quality checks
check: fmt vet lint test-race

# Install the binary
install: build
	cp $(BINARY_NAME) /usr/local/bin/

# Uninstall the binary
uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  all         - Clean, get dependencies, test, and build"
	@echo "  build       - Build the binary"
	@echo "  build-windows - Build Windows binary"
	@echo "  build-linux - Build Linux binary"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race detector"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
	@echo "  vet         - Run go vet"
	@echo "  install-dev - Install development dependencies"
	@echo "  check       - Run all quality checks"
	@echo "  install     - Install binary to /usr/local/bin"
	@echo "  uninstall   - Remove binary from /usr/local/bin"
	@echo "  help        - Show this help message"