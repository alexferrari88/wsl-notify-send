@echo off
echo Running wsl-notify-send tests...
echo.

REM Run tests with coverage
echo Running tests with coverage...
go test -v -coverprofile=coverage.out ./...

REM Check if tests passed
if %errorlevel% neq 0 (
    echo.
    echo Tests failed!
    exit /b 1
)

REM Generate coverage report
echo.
echo Generating coverage report...
go tool cover -html=coverage.out -o coverage.html

REM Show coverage summary
echo.
echo Coverage summary:
go tool cover -func=coverage.out

echo.
echo Tests completed successfully!
echo Coverage report generated: coverage.html