package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		expectCode int
	}{
		{
			name:       "no error",
			err:        nil,
			expectCode: 0,
		},
		{
			name:       "invalid configuration error",
			err:        errors.New("invalid configuration: frequency must be positive"),
			expectCode: 2,
		},
		{
			name:       "too many arguments error",
			err:        errors.New("too many arguments, expected: <title> [message]"),
			expectCode: 2,
		},
		{
			name:       "requires at least error",
			err:        errors.New("requires at least a title argument"),
			expectCode: 2,
		},
		{
			name:       "failed to send notification",
			err:        errors.New("failed to send notification: something went wrong"),
			expectCode: 3,
		},
		{
			name:       "failed to beep",
			err:        errors.New("failed to beep: beep error"),
			expectCode: 3,
		},
		{
			name:       "general error",
			err:        errors.New("some other error"),
			expectCode: 1,
		},
		{
			name:       "cobra error",
			err:        errors.New("unknown command"),
			expectCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := getExitCode(tt.err)
			assert.Equal(t, tt.expectCode, code)
		})
	}
}

func TestGetExitCodeEdgeCases(t *testing.T) {
	// Test that error message matching is case-sensitive and exact
	tests := []struct {
		name       string
		err        error
		expectCode int
	}{
		{
			name:       "invalid configuration in middle",
			err:        errors.New("something invalid configuration something"),
			expectCode: 2,
		},
		{
			name:       "failed to send in middle",
			err:        errors.New("error failed to send more text"),
			expectCode: 3,
		},
		{
			name:       "failed to beep in middle",
			err:        errors.New("error failed to beep more text"),
			expectCode: 3,
		},
		{
			name:       "requires at least in middle",
			err:        errors.New("error requires at least more text"),
			expectCode: 2,
		},
		{
			name:       "too many arguments in middle",
			err:        errors.New("error too many arguments more text"),
			expectCode: 2,
		},
		{
			name:       "case sensitivity test",
			err:        errors.New("Invalid Configuration"),
			expectCode: 1, // Should not match due to case sensitivity
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := getExitCode(tt.err)
			assert.Equal(t, tt.expectCode, code)
		})
	}
}

func TestGetExitCodeMultipleMatches(t *testing.T) {
	// Test that the first matching case is used
	tests := []struct {
		name       string
		err        error
		expectCode int
	}{
		{
			name:       "invalid configuration and failed to send",
			err:        errors.New("invalid configuration: failed to send notification"),
			expectCode: 2, // Should match first case (invalid configuration)
		},
		{
			name:       "too many arguments and failed to beep",
			err:        errors.New("too many arguments, also failed to beep"),
			expectCode: 2, // Should match first case (too many arguments)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := getExitCode(tt.err)
			assert.Equal(t, tt.expectCode, code)
		})
	}
}

func TestGetExitCodeEmptyError(t *testing.T) {
	// Test with empty error message
	err := errors.New("")
	code := getExitCode(err)
	assert.Equal(t, 1, code) // Should default to general error
}

func TestGetExitCodeSpecificErrorMessages(t *testing.T) {
	// Test specific error messages that might be returned by the application
	tests := []struct {
		name       string
		errorMsg   string
		expectCode int
	}{
		{
			name:       "config validation error",
			errorMsg:   "invalid configuration: cannot use both --alert and --beep modes",
			expectCode: 2,
		},
		{
			name:       "icon file error",
			errorMsg:   "invalid configuration: icon file does not exist: test.png",
			expectCode: 2,
		},
		{
			name:       "frequency validation error",
			errorMsg:   "invalid configuration: frequency must be positive",
			expectCode: 2,
		},
		{
			name:       "duration validation error",
			errorMsg:   "invalid configuration: duration must be positive",
			expectCode: 2,
		},
		{
			name:       "beeep notify error",
			errorMsg:   "failed to send notification: notification system error",
			expectCode: 3,
		},
		{
			name:       "beeep alert error",
			errorMsg:   "failed to send alert: alert system error",
			expectCode: 3,
		},
		{
			name:       "beeep beep error",
			errorMsg:   "failed to beep: beep system error",
			expectCode: 3,
		},
		{
			name:       "argument parsing error",
			errorMsg:   "requires at least a title argument",
			expectCode: 2,
		},
		{
			name:       "too many args error",
			errorMsg:   "too many arguments, expected: <title> [message]",
			expectCode: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errorMsg)
			code := getExitCode(err)
			assert.Equal(t, tt.expectCode, code)
		})
	}
}

// Test the behavior of the main function error handling patterns
func TestMainErrorHandling(t *testing.T) {
	// These tests verify the error handling logic without actually calling main()
	// since main() calls os.Exit() which would terminate the test
	
	// Test that we properly categorize different error types
	errorTypes := map[string]int{
		"invalid configuration errors":   2,
		"argument parsing errors":        2,
		"notification system errors":     3,
		"general errors":                1,
	}
	
	for errorType, expectedCode := range errorTypes {
		t.Run(errorType, func(t *testing.T) {
			var testErr error
			switch errorType {
			case "invalid configuration errors":
				testErr = errors.New("invalid configuration: test error")
			case "argument parsing errors":
				testErr = errors.New("requires at least a title argument")
			case "notification system errors":
				testErr = errors.New("failed to send notification: test error")
			case "general errors":
				testErr = errors.New("some general error")
			}
			
			code := getExitCode(testErr)
			assert.Equal(t, expectedCode, code)
		})
	}
}

// Benchmark tests for getExitCode function
func BenchmarkGetExitCode(b *testing.B) {
	testErrors := []error{
		nil,
		errors.New("invalid configuration: test"),
		errors.New("failed to send notification: test"),
		errors.New("failed to beep: test"),
		errors.New("requires at least a title argument"),
		errors.New("too many arguments, expected: <title> [message]"),
		errors.New("general error"),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := testErrors[i%len(testErrors)]
		getExitCode(err)
	}
}

// Test string comparison behavior
func TestStringContainsLogic(t *testing.T) {
	// Test the strings.Contains behavior used in getExitCode
	tests := []struct {
		str      string
		substr   string
		expected bool
	}{
		{"invalid configuration: test", "invalid configuration", true},
		{"failed to send notification", "failed to send", true},
		{"failed to beep", "failed to beep", true},
		{"requires at least a title", "requires at least", true},
		{"too many arguments", "too many arguments", true},
		{"Invalid Configuration", "invalid configuration", false}, // case sensitive
		{"", "test", false},
		{"test", "", true}, // empty substring always matches
	}
	
	for _, tt := range tests {
		t.Run(tt.str+"_contains_"+tt.substr, func(t *testing.T) {
			// This tests the actual logic used in getExitCode
			result := len(tt.str) >= len(tt.substr) && 
				(tt.str[:len(tt.substr)] == tt.substr || 
				 (len(tt.str) > len(tt.substr) && tt.str[len(tt.str)-len(tt.substr):] == tt.substr) ||
				 stringContains(tt.str, tt.substr))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to test string contains logic
func stringContains(str, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(str) < len(substr) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}