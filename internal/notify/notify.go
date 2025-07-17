package notify

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gen2brain/beeep"
)

// Default beeper instance
var defaultBeeper Beeper = NewDefaultBeeper()

// SetBeeper allows setting a custom beeper (mainly for testing)
func SetBeeper(b Beeper) {
	defaultBeeper = b
}

// Notify sends a desktop notification without sound
func Notify(title, message, icon, appName string) error {
	// Set application name if provided
	if appName != "" {
		defaultBeeper.SetAppName(appName)
	}

	// Process icon
	iconData, err := processIcon(icon)
	if err != nil {
		return fmt.Errorf("failed to process icon: %w", err)
	}

	// Send notification
	if err := defaultBeeper.Notify(title, message, iconData); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

// Alert sends a desktop notification with sound
func Alert(title, message, icon, appName string) error {
	// Set application name if provided
	if appName != "" {
		defaultBeeper.SetAppName(appName)
	}

	// Process icon
	iconData, err := processIcon(icon)
	if err != nil {
		return fmt.Errorf("failed to process icon: %w", err)
	}

	// Send alert
	if err := defaultBeeper.Alert(title, message, iconData); err != nil {
		return fmt.Errorf("failed to send alert: %w", err)
	}

	return nil
}

// Beep plays a beep sound
func Beep(frequency float64, duration int) error {
	if err := defaultBeeper.Beep(frequency, duration); err != nil {
		return fmt.Errorf("failed to beep: %w", err)
	}

	return nil
}

// Wrapper functions for the actual beeep library
func beepNotify(title, message string, icon interface{}) error {
	return beeep.Notify(title, message, icon)
}

func beepAlert(title, message string, icon interface{}) error {
	return beeep.Alert(title, message, icon)
}

func beepBeep(freq float64, duration int) error {
	return beeep.Beep(freq, duration)
}

func beepSetAppName(name string) {
	beeep.AppName = name
}

// processIcon handles icon processing for notifications
func processIcon(icon string) (interface{}, error) {
	// No icon specified
	if icon == "" {
		return "", nil
	}

	// Check if it's an absolute path
	if filepath.IsAbs(icon) {
		// Read file data
		data, err := os.ReadFile(icon)
		if err != nil {
			return nil, fmt.Errorf("cannot read icon file: %w", err)
		}
		return data, nil
	}

	// Check if it's a relative path with directory separators
	if filepath.Dir(icon) != "." {
		// Read file data
		data, err := os.ReadFile(icon)
		if err != nil {
			return nil, fmt.Errorf("cannot read icon file: %w", err)
		}
		return data, nil
	}

	// Check if it has a file extension and exists as a file
	if filepath.Ext(icon) != "" {
		if _, err := os.Stat(icon); err == nil {
			// File exists, read it
			data, err := os.ReadFile(icon)
			if err != nil {
				return nil, fmt.Errorf("cannot read icon file: %w", err)
			}
			return data, nil
		}
	}

	// Treat as stock icon name
	return icon, nil
}

// GetSupportedFormats returns the supported icon formats
func GetSupportedFormats() []string {
	return []string{".png", ".jpg", ".jpeg", ".ico", ".bmp"}
}

// IsValidIconFormat checks if the given file extension is supported
func IsValidIconFormat(ext string) bool {
	supported := GetSupportedFormats()
	for _, format := range supported {
		if ext == format {
			return true
		}
	}
	return false
}
