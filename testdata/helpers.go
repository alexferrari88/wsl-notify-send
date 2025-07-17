package testdata

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTestIcon creates a temporary icon file for testing
func CreateTestIcon(t *testing.T, ext string) string {
	tempDir := t.TempDir()
	iconPath := filepath.Join(tempDir, "test"+ext)
	
	content := "dummy_" + ext[1:] + "_content"
	err := os.WriteFile(iconPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test icon: %v", err)
	}
	
	return iconPath
}

// CreateTestIcons creates multiple test icon files
func CreateTestIcons(t *testing.T, extensions []string) map[string]string {
	icons := make(map[string]string)
	
	for _, ext := range extensions {
		icons[ext] = CreateTestIcon(t, ext)
	}
	
	return icons
}

// SupportedIconFormats returns all supported icon formats
func SupportedIconFormats() []string {
	return []string{".png", ".jpg", ".jpeg", ".ico", ".bmp"}
}

// UnsupportedIconFormats returns unsupported icon formats for testing
func UnsupportedIconFormats() []string {
	return []string{".gif", ".webp", ".svg", ".txt", ".exe"}
}

// CreateTempFileWithContent creates a temporary file with given content
func CreateTempFileWithContent(t *testing.T, filename, content string) string {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, filename)
	
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	return filePath
}

// ValidConfigDefaults returns a valid default configuration for testing
func ValidConfigDefaults() map[string]interface{} {
	return map[string]interface{}{
		"Frequency": 587.0,
		"Duration":  500,
		"AppName":   "wsl-notify-send",
		"Quiet":     false,
		"AlertMode": false,
		"BeepMode":  false,
		"Icon":      "",
		"Version":   false,
	}
}

// CommonErrorMessages returns common error messages for testing
func CommonErrorMessages() map[string]string {
	return map[string]string{
		"InvalidConfig":      "invalid configuration",
		"TooManyArgs":        "too many arguments",
		"RequiresTitle":      "requires at least a title argument",
		"NotificationFailed": "failed to send notification",
		"AlertFailed":        "failed to send alert",
		"BeepFailed":         "failed to beep",
		"IconNotFound":       "icon file does not exist",
		"UnsupportedFormat":  "unsupported icon format",
		"FrequencyPositive":  "frequency must be positive",
		"DurationPositive":   "duration must be positive",
		"BothModes":          "cannot use both --alert and --beep modes",
	}
}

// TestCLIArgs returns common CLI argument combinations for testing
func TestCLIArgs() map[string][]string {
	return map[string][]string{
		"basic":           {"Test", "Message"},
		"titleOnly":       {"Just Title"},
		"alert":           {"--alert", "Alert", "Message"},
		"beep":            {"--beep"},
		"withIcon":        {"--icon", "warning", "Test", "Message"},
		"withAppName":     {"--app-name", "MyApp", "Test", "Message"},
		"quiet":           {"--quiet", "Test", "Message"},
		"version":         {"--version"},
		"help":            {"--help"},
		"customBeep":      {"--beep", "--freq", "1000", "--duration", "1000"},
		"shortFlags":      {"-a", "-i", "warning", "-q", "Test", "Message"},
		"noArgs":          {},
		"tooManyArgs":     {"title", "message", "extra", "args"},
		"invalidCombo":    {"--alert", "--beep"},
	}
}

// ExpectedExitCodes returns expected exit codes for different scenarios
func ExpectedExitCodes() map[string]int {
	return map[string]int{
		"success":              0,
		"general_error":        1,
		"invalid_args":         2,
		"notification_failed":  3,
	}
}