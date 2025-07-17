package main

import (
	"fmt"
	"os"
	"strings"
	"wsl-notify-send/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exitCode := getExitCode(err)

		// Only print error if not in quiet mode
		if !cmd.IsQuietMode() {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		os.Exit(exitCode)
	}
}

func getExitCode(err error) int {
	if err == nil {
		return 0
	}

	errStr := err.Error()

	// Check for different error types
	switch {
	case strings.Contains(errStr, "invalid configuration") || strings.Contains(errStr, "too many arguments") || strings.Contains(errStr, "requires at least"):
		return 2 // Invalid arguments
	case strings.Contains(errStr, "failed to send") || strings.Contains(errStr, "failed to beep"):
		return 3 // Notification failed
	default:
		return 1 // General error
	}
}
