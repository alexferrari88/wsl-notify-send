package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"wsl-notify-send/internal/config"
	"wsl-notify-send/internal/notify"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBeeper for CLI tests
type MockBeeper struct {
	mock.Mock
}

func (m *MockBeeper) Notify(title, message string, icon interface{}) error {
	args := m.Called(title, message, icon)
	return args.Error(0)
}

func (m *MockBeeper) Alert(title, message string, icon interface{}) error {
	args := m.Called(title, message, icon)
	return args.Error(0)
}

func (m *MockBeeper) Beep(freq float64, duration int) error {
	args := m.Called(freq, duration)
	return args.Error(0)
}

func (m *MockBeeper) SetAppName(name string) {
	m.Called(name)
}

func setupMockBeeper(t *testing.T) *MockBeeper {
	mockBeeper := new(MockBeeper)
	notify.SetBeeper(mockBeeper)

	// Initialize config with default values, resetting all fields
	cfg = config.Config{
		AlertMode: false,
		BeepMode:  false,
		Icon:      "",
		AppName:   "wsl-notify-send",
		Frequency: 587.0,
		Duration:  500,
		Quiet:     false,
		Version:   false,
	}

	// Reset to default beeper after test
	t.Cleanup(func() {
		notify.SetBeeper(notify.NewDefaultBeeper())
		// Reset global config with default values
		cfg = config.Config{
			AlertMode: false,
			BeepMode:  false,
			Icon:      "",
			AppName:   "wsl-notify-send",
			Frequency: 587.0,
			Duration:  500,
			Quiet:     false,
			Version:   false,
		}
	})

	return mockBeeper
}

// Helper function to execute command and capture output
func executeCommand(args []string) (string, error) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set test args
	os.Args = append([]string{"wsl-notify-send"}, args...)

	// Capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Set args for command
	rootCmd.SetArgs(args)

	// Execute command
	err := rootCmd.Execute()

	// Capture the output and config state before resetting
	outputStr := buf.String()

	// Reset command for next test
	rootCmd.SetArgs(nil)
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	// Reset command flags to defaults
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})

	// Also reset the global config to match the default flags
	cfg.AlertMode = false
	cfg.BeepMode = false
	cfg.Icon = ""
	cfg.AppName = "wsl-notify-send"
	cfg.Frequency = 587.0
	cfg.Duration = 500
	cfg.Quiet = false
	cfg.Version = false

	return outputStr, err
}

func TestRootCommand_BasicNotification(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Notify", "Test Title", "Test Message", "").Return(nil).Once()

	_, err := executeCommand([]string{"Test Title", "Test Message"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_NotificationWithAppName(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "MyApp").Once()
	mockBeeper.On("Notify", "Test Title", "Test Message", "").Return(nil).Once()

	_, err := executeCommand([]string{"--app-name", "MyApp", "Test Title", "Test Message"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_AlertMode(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Alert", "Alert Title", "Alert Message", "").Return(nil).Once()

	_, err := executeCommand([]string{"--alert", "Alert Title", "Alert Message"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_BeepMode(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("Beep", 587.0, 500).Return(nil).Once()

	_, err := executeCommand([]string{"--beep"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_BeepModeWithCustomParams(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("Beep", 1000.0, 1000).Return(nil).Once()

	_, err := executeCommand([]string{"--beep", "--freq", "1000", "--duration", "1000"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_TitleOnly(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Notify", "Just Title", "", "").Return(nil).Once()

	_, err := executeCommand([]string{"Just Title"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_WithIcon(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Notify", "Title", "Message", "warning").Return(nil).Once()

	_, err := executeCommand([]string{"--icon", "warning", "Title", "Message"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_VersionFlag(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	output, err := executeCommand([]string{"--version"})

	assert.NoError(t, err)
	assert.Contains(t, output, "wsl-notify-send version "+Version)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_HelpFlag(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	output, err := executeCommand([]string{"--help"})

	assert.NoError(t, err)
	assert.Contains(t, output, "wsl-notify-send is a cross-platform notification tool")
	assert.Contains(t, output, "Usage:")
	assert.Contains(t, output, "Examples:")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_NoArgumentsError(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	_, err := executeCommand([]string{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires at least a title argument")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_TooManyArgumentsError(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	_, err := executeCommand([]string{"title", "message", "extra"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too many arguments")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_InvalidConfiguration(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	_, err := executeCommand([]string{"--alert", "--beep"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
	assert.Contains(t, err.Error(), "cannot use both --alert and --beep modes")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_NotificationFailure(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Notify", "Test", "Message", "").Return(errors.New("notification failed")).Once()

	_, err := executeCommand([]string{"Test", "Message"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send notification")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_BeepFailure(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("Beep", 587.0, 500).Return(errors.New("beep failed")).Once()

	_, err := executeCommand([]string{"--beep"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to beep")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_IconFileValidation(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.png")

	_, err := executeCommand([]string{"--icon", nonExistentFile, "Test", "Message"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
	assert.Contains(t, err.Error(), "icon file does not exist")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_FrequencyValidation(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	_, err := executeCommand([]string{"--beep", "--freq", "-100"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
	assert.Contains(t, err.Error(), "frequency must be positive")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_DurationValidation(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	_, err := executeCommand([]string{"--beep", "--duration", "-100"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
	assert.Contains(t, err.Error(), "duration must be positive")
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_QuietModeFlag(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Notify", "Test", "Message", "").Return(nil).Once()

	// Execute command and capture config state before it gets reset
	rootCmd.SetArgs([]string{"--quiet", "Test", "Message"})
	err := rootCmd.Execute()

	// Check the config state immediately after execution
	assert.NoError(t, err)
	assert.True(t, cfg.Quiet)
	mockBeeper.AssertExpectations(t)

	// Reset for next test
	rootCmd.SetArgs(nil)
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})
	cfg.AlertMode = false
	cfg.BeepMode = false
	cfg.Icon = ""
	cfg.AppName = "wsl-notify-send"
	cfg.Frequency = 587.0
	cfg.Duration = 500
	cfg.Quiet = false
	cfg.Version = false
}

func TestRootCommand_ShortFlags(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Alert", "Alert", "Message", "warning").Return(nil).Once()

	// Execute command and capture config state before it gets reset
	rootCmd.SetArgs([]string{"-a", "-i", "warning", "-q", "Alert", "Message"})
	err := rootCmd.Execute()

	// Check the config state immediately after execution
	assert.NoError(t, err)
	assert.True(t, cfg.AlertMode)
	assert.True(t, cfg.Quiet)
	assert.Equal(t, "warning", cfg.Icon)
	mockBeeper.AssertExpectations(t)

	// Reset for next test
	rootCmd.SetArgs(nil)
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})
	cfg.AlertMode = false
	cfg.BeepMode = false
	cfg.Icon = ""
	cfg.AppName = "wsl-notify-send"
	cfg.Frequency = 587.0
	cfg.Duration = 500
	cfg.Quiet = false
	cfg.Version = false
}

func TestRootCommand_BeepModeWithTitle(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("Beep", 587.0, 500).Return(nil).Once()

	// In beep mode, title arguments should be ignored
	_, err := executeCommand([]string{"--beep", "This", "Should", "Be", "Ignored"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_VersionModeWithTitle(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	// In version mode, title arguments should be ignored
	output, err := executeCommand([]string{"--version", "This", "Should", "Be", "Ignored"})

	assert.NoError(t, err)
	assert.Contains(t, output, "wsl-notify-send version "+Version)
	mockBeeper.AssertExpectations(t)
}

func TestIsQuietMode(t *testing.T) {
	// Test default state
	assert.False(t, IsQuietMode())

	// Test when quiet mode is enabled
	cfg.Quiet = true
	assert.True(t, IsQuietMode())

	// Reset
	cfg.Quiet = false
	assert.False(t, IsQuietMode())
}

func TestRootCommand_DefaultAppName(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	mockBeeper.On("SetAppName", "wsl-notify-send").Once()
	mockBeeper.On("Notify", "Test", "Message", "").Return(nil).Once()

	_, err := executeCommand([]string{"--app-name", "wsl-notify-send", "Test", "Message"})

	assert.NoError(t, err)
	mockBeeper.AssertExpectations(t)
}

func TestRootCommand_AllFlagsTogethert(t *testing.T) {
	mockBeeper := setupMockBeeper(t)

	tempDir := t.TempDir()
	iconFile := filepath.Join(tempDir, "test.png")
	err := os.WriteFile(iconFile, []byte("test content"), 0644)
	require.NoError(t, err)

	iconContent := []byte("test content")

	mockBeeper.On("SetAppName", "TestApp").Once()
	mockBeeper.On("Alert", "Test", "Message", iconContent).Return(nil).Once()

	// Execute command and capture config state before it gets reset
	rootCmd.SetArgs([]string{
		"--alert",
		"--icon", iconFile,
		"--app-name", "TestApp",
		"--quiet",
		"Test", "Message",
	})
	err = rootCmd.Execute()

	// Check the config state immediately after execution
	assert.NoError(t, err)
	assert.True(t, cfg.AlertMode)
	assert.True(t, cfg.Quiet)
	assert.Equal(t, "TestApp", cfg.AppName)
	assert.Equal(t, iconFile, cfg.Icon)
	mockBeeper.AssertExpectations(t)

	// Reset for next test
	rootCmd.SetArgs(nil)
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})
	cfg.AlertMode = false
	cfg.BeepMode = false
	cfg.Icon = ""
	cfg.AppName = "wsl-notify-send"
	cfg.Frequency = 587.0
	cfg.Duration = 500
	cfg.Quiet = false
	cfg.Version = false
}
