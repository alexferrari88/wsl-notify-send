# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

wsl-notify-send is a cross-platform notification tool for Windows and WSL2 built in Go. It provides desktop notifications, alerts, and beeps using the beeep library with a clean CLI interface built on Cobra.

## Development Commands

### Building
- `go build -o wsl-notify-send.exe .` - Build for current platform
- `make build` - Build using Makefile
- `make build-windows` - Cross-compile for Windows
- `make build-linux` - Cross-compile for Linux

### Testing
- `go test ./...` - Run all tests
- `make test` - Run tests via Makefile
- `make test-race` - Run tests with race detector
- `make coverage` - Generate coverage report (creates coverage.out and coverage.html)

### Code Quality
- `go fmt ./...` - Format code
- `make fmt` - Format via Makefile
- `go vet ./...` - Run go vet
- `make vet` - Vet via Makefile
- `golangci-lint run` - Run linter (requires golangci-lint)
- `make lint` - Run linter via Makefile
- `make check` - Run all quality checks (fmt, vet, lint, test-race)

### Dependencies
- `go mod download && go mod tidy` - Download and tidy dependencies
- `make deps` - Dependencies via Makefile
- `make install-dev` - Install development dependencies (golangci-lint)

### Utilities
- `make clean` - Clean build artifacts
- `make all` - Full build pipeline (clean, deps, test, build)

## Architecture

### Core Components
- **main.go** - Entry point with error handling and exit code logic
- **cmd/root.go** - Cobra CLI command definition with flag handling
- **internal/config/config.go** - Configuration struct and validation logic
- **internal/notify/notify.go** - Core notification functions wrapping beeep library
- **internal/notify/interface.go** - Beeper interface for testability

### Key Patterns
- Uses dependency injection with Beeper interface for testing
- Error handling with specific exit codes (0=success, 1=general, 2=invalid args, 3=notification failed)
- Icon processing supports file paths, stock icons, and embedded data
- Quiet mode suppresses error output while preserving exit codes

### Test Structure
- All packages have corresponding _test.go files
- Uses testify for assertions and mocking
- testdata/ directory contains test assets (icons, files)
- test.sh and test.bat provide cross-platform test runners

### Exit Codes
- 0: Success
- 1: General error
- 2: Invalid arguments or configuration
- 3: Notification failed to send