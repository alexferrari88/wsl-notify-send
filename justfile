# Go parameters
gocmd := "go"
gobuild := gocmd + " build"
goclean := gocmd + " clean"
gotest := gocmd + " test"
goget := gocmd + " get"
gomod := gocmd + " mod"
binary_name := "wsl-notify-send"
binary_windows := binary_name + ".exe"

# Show available recipes
default:
    @just --list

# Clean, get dependencies, test, and build
all: clean deps test build

# Build the binary
build:
    {{gobuild}} -o {{binary_name}} -v .

# Build Windows binary
build-windows:
    env GOOS=windows GOARCH=amd64 {{gobuild}} -o {{binary_windows}} -v .

# Build Linux binary
build-linux:
    env GOOS=linux GOARCH=amd64 {{gobuild}} -o {{binary_name}} -v .

# Clean build artifacts
clean:
    {{goclean}}
    rm -f {{binary_name}}
    rm -f {{binary_windows}}
    rm -f coverage.out
    rm -f coverage.html

# Run tests
test:
    {{gotest}} -v ./...

# Run tests with race detector
test-race:
    {{gotest}} -v -race ./...

# Run tests with coverage report
coverage:
    {{gotest}} -v -coverprofile=coverage.out ./...
    {{gocmd}} tool cover -html=coverage.out -o coverage.html
    {{gocmd}} tool cover -func=coverage.out

# Download and tidy dependencies
deps:
    {{gomod}} download
    {{gomod}} tidy

# Format code
fmt:
    {{gocmd}} fmt ./...

# Run linter
lint:
    golangci-lint run

# Run go vet
vet:
    {{gocmd}} vet ./...

# Install development dependencies
install-dev:
    {{goget}} github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all quality checks
check: fmt vet lint test-race

# Install binary to /usr/local/bin
install: build
    cp {{binary_name}} /usr/local/bin/

# Remove binary from /usr/local/bin
uninstall:
    rm -f /usr/local/bin/{{binary_name}}