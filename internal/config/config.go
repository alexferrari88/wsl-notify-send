package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	// Mode flags
	AlertMode bool
	BeepMode  bool

	// Content options
	Icon    string
	AppName string

	// Beep options
	Frequency float64
	Duration  int

	// Utility options
	Quiet   bool
	Version bool
}

func (c *Config) Validate() error {
	// Can't have both alert and beep mode
	if c.AlertMode && c.BeepMode {
		return errors.New("cannot use both --alert and --beep modes")
	}

	// Validate icon file if provided
	if c.Icon != "" {
		if err := c.validateIcon(); err != nil {
			return err
		}
	}

	// Validate beep parameters
	if c.Frequency <= 0 {
		return errors.New("frequency must be positive")
	}

	if c.Duration <= 0 {
		return errors.New("duration must be positive")
	}

	return nil
}

func (c *Config) validateIcon() error {
	// Check if it's an absolute path
	if filepath.IsAbs(c.Icon) {
		return c.validateIconFile()
	}

	// Check if it's a relative path with directory separators
	if filepath.Dir(c.Icon) != "." {
		return c.validateIconFile()
	}

	// Check if it has a file extension and exists as a file
	if filepath.Ext(c.Icon) != "" {
		if _, err := os.Stat(c.Icon); err == nil {
			return c.validateIconFile()
		}
	}

	// Treat as stock icon name - no validation needed
	return nil
}

func (c *Config) validateIconFile() error {
	// Check if file exists
	if _, err := os.Stat(c.Icon); err != nil {
		if os.IsNotExist(err) {
			return errors.New("icon file does not exist: " + c.Icon)
		}
		return errors.New("cannot access icon file: " + err.Error())
	}

	// Check file extension
	ext := filepath.Ext(c.Icon)
	switch ext {
	case ".png", ".jpg", ".jpeg", ".ico", ".bmp":
		return nil
	default:
		return errors.New("unsupported icon format: " + ext + " (supported: .png, .jpg, .jpeg, .ico, .bmp)")
	}
}
