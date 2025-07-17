#!/bin/bash

echo "Running wsl-notify-send tests..."
echo

# Run tests with coverage
echo "Running tests with coverage..."
go test -v -coverprofile=coverage.out ./...

# Check if tests passed
if [ $? -ne 0 ]; then
    echo
    echo "Tests failed!"
    exit 1
fi

# Generate coverage report
echo
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Show coverage summary
echo
echo "Coverage summary:"
go tool cover -func=coverage.out

echo
echo "Tests completed successfully!"
echo "Coverage report generated: coverage.html"