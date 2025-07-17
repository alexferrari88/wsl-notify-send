package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid default config",
			config: Config{
				Frequency: 587.0,
				Duration:  500,
				AppName:   "test-app",
			},
			expectError: false,
		},
		{
			name: "both alert and beep mode",
			config: Config{
				AlertMode: true,
				BeepMode:  true,
				Frequency: 587.0,
				Duration:  500,
			},
			expectError: true,
			errorMsg:    "cannot use both --alert and --beep modes",
		},
		{
			name: "negative frequency",
			config: Config{
				Frequency: -100,
				Duration:  500,
			},
			expectError: true,
			errorMsg:    "frequency must be positive",
		},
		{
			name: "zero frequency",
			config: Config{
				Frequency: 0,
				Duration:  500,
			},
			expectError: true,
			errorMsg:    "frequency must be positive",
		},
		{
			name: "negative duration",
			config: Config{
				Frequency: 587.0,
				Duration:  -100,
			},
			expectError: true,
			errorMsg:    "duration must be positive",
		},
		{
			name: "zero duration",
			config: Config{
				Frequency: 587.0,
				Duration:  0,
			},
			expectError: true,
			errorMsg:    "duration must be positive",
		},
		{
			name: "valid alert mode",
			config: Config{
				AlertMode: true,
				Frequency: 587.0,
				Duration:  500,
			},
			expectError: false,
		},
		{
			name: "valid beep mode",
			config: Config{
				BeepMode:  true,
				Frequency: 587.0,
				Duration:  500,
			},
			expectError: false,
		},
		{
			name: "stock icon name",
			config: Config{
				Icon:      "warning",
				Frequency: 587.0,
				Duration:  500,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_ValidateIcon(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()

	// Create valid PNG file
	validPngPath := filepath.Join(tempDir, "test.png")
	err := os.WriteFile(validPngPath, []byte("dummy png content"), 0644)
	require.NoError(t, err)

	// Create invalid format file
	invalidFormatPath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(invalidFormatPath, []byte("text content"), 0644)
	require.NoError(t, err)

	// Non-existent file path
	nonExistentPath := filepath.Join(tempDir, "nonexistent.png")

	tests := []struct {
		name        string
		iconPath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty icon path",
			iconPath:    "",
			expectError: false,
		},
		{
			name:        "stock icon name",
			iconPath:    "warning",
			expectError: false,
		},
		{
			name:        "valid PNG file",
			iconPath:    validPngPath,
			expectError: false,
		},
		{
			name:        "non-existent file",
			iconPath:    nonExistentPath,
			expectError: true,
			errorMsg:    "icon file does not exist",
		},
		{
			name:        "unsupported format",
			iconPath:    invalidFormatPath,
			expectError: true,
			errorMsg:    "unsupported icon format: .txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Icon:      tt.iconPath,
				Frequency: 587.0,
				Duration:  500,
			}

			err := config.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_ValidateIconFormats(t *testing.T) {
	tempDir := t.TempDir()

	supportedFormats := []string{".png", ".jpg", ".jpeg", ".ico", ".bmp"}

	for _, format := range supportedFormats {
		t.Run("supported format "+format, func(t *testing.T) {
			filePath := filepath.Join(tempDir, "test"+format)
			err := os.WriteFile(filePath, []byte("dummy content"), 0644)
			require.NoError(t, err)

			config := Config{
				Icon:      filePath,
				Frequency: 587.0,
				Duration:  500,
			}

			err = config.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestConfig_ValidateIconFormatsCase(t *testing.T) {
	tempDir := t.TempDir()

	// Test case sensitivity
	upperCaseFile := filepath.Join(tempDir, "test.PNG")
	err := os.WriteFile(upperCaseFile, []byte("dummy content"), 0644)
	require.NoError(t, err)

	config := Config{
		Icon:      upperCaseFile,
		Frequency: 587.0,
		Duration:  500,
	}

	err = config.Validate()
	// This should fail because the validation is case-sensitive
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported icon format: .PNG")
}

func TestConfig_ValidateFrequencyBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		frequency float64
		valid     bool
	}{
		{"very small positive", 0.001, true},
		{"exactly zero", 0.0, false},
		{"small negative", -0.001, false},
		{"large positive", 20000.0, true},
		{"normal frequency", 587.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Frequency: tt.frequency,
				Duration:  500,
			}

			err := config.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "frequency must be positive")
			}
		})
	}
}

func TestConfig_ValidateDurationBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		valid    bool
	}{
		{"minimum valid", 1, true},
		{"exactly zero", 0, false},
		{"negative", -1, false},
		{"large positive", 10000, true},
		{"normal duration", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Frequency: 587.0,
				Duration:  tt.duration,
			}

			err := config.Validate()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "duration must be positive")
			}
		})
	}
}
