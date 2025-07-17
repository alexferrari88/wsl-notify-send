package notify

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBeeper is a mock implementation of the Beeper interface
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
	SetBeeper(mockBeeper)
	
	// Reset to default beeper after test
	t.Cleanup(func() {
		SetBeeper(NewDefaultBeeper())
	})
	
	return mockBeeper
}

func TestNotify(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		message     string
		icon        string
		appName     string
		mockSetup   func(*MockBeeper)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful notification",
			title:   "Test Title",
			message: "Test Message",
			icon:    "",
			appName: "TestApp",
			mockSetup: func(m *MockBeeper) {
				m.On("SetAppName", "TestApp").Once()
				m.On("Notify", "Test Title", "Test Message", "").Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:    "successful notification without app name",
			title:   "Test Title",
			message: "Test Message",
			icon:    "",
			appName: "",
			mockSetup: func(m *MockBeeper) {
				m.On("Notify", "Test Title", "Test Message", "").Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:    "notification failure",
			title:   "Test Title",
			message: "Test Message",
			icon:    "",
			appName: "",
			mockSetup: func(m *MockBeeper) {
				m.On("Notify", "Test Title", "Test Message", "").Return(errors.New("notification failed")).Once()
			},
			expectError: true,
			errorMsg:    "failed to send notification",
		},
		{
			name:    "stock icon notification",
			title:   "Test Title",
			message: "Test Message",
			icon:    "warning",
			appName: "",
			mockSetup: func(m *MockBeeper) {
				m.On("Notify", "Test Title", "Test Message", "warning").Return(nil).Once()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBeeper := setupMockBeeper(t)
			tt.mockSetup(mockBeeper)

			err := Notify(tt.title, tt.message, tt.icon, tt.appName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockBeeper.AssertExpectations(t)
		})
	}
}

func TestAlert(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		message     string
		icon        string
		appName     string
		mockSetup   func(*MockBeeper)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful alert",
			title:   "Alert Title",
			message: "Alert Message",
			icon:    "",
			appName: "TestApp",
			mockSetup: func(m *MockBeeper) {
				m.On("SetAppName", "TestApp").Once()
				m.On("Alert", "Alert Title", "Alert Message", "").Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:    "alert failure",
			title:   "Alert Title",
			message: "Alert Message",
			icon:    "",
			appName: "",
			mockSetup: func(m *MockBeeper) {
				m.On("Alert", "Alert Title", "Alert Message", "").Return(errors.New("alert failed")).Once()
			},
			expectError: true,
			errorMsg:    "failed to send alert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBeeper := setupMockBeeper(t)
			tt.mockSetup(mockBeeper)

			err := Alert(tt.title, tt.message, tt.icon, tt.appName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockBeeper.AssertExpectations(t)
		})
	}
}

func TestBeep(t *testing.T) {
	tests := []struct {
		name        string
		frequency   float64
		duration    int
		mockSetup   func(*MockBeeper)
		expectError bool
		errorMsg    string
	}{
		{
			name:      "successful beep",
			frequency: 587.0,
			duration:  500,
			mockSetup: func(m *MockBeeper) {
				m.On("Beep", 587.0, 500).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:      "beep failure",
			frequency: 587.0,
			duration:  500,
			mockSetup: func(m *MockBeeper) {
				m.On("Beep", 587.0, 500).Return(errors.New("beep failed")).Once()
			},
			expectError: true,
			errorMsg:    "failed to beep",
		},
		{
			name:      "custom frequency and duration",
			frequency: 1000.0,
			duration:  1000,
			mockSetup: func(m *MockBeeper) {
				m.On("Beep", 1000.0, 1000).Return(nil).Once()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBeeper := setupMockBeeper(t)
			tt.mockSetup(mockBeeper)

			err := Beep(tt.frequency, tt.duration)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockBeeper.AssertExpectations(t)
		})
	}
}

func TestProcessIcon(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test file
	testFile := filepath.Join(tempDir, "test.png")
	testContent := []byte("test icon content")
	err := os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)
	
	// Non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.png")

	tests := []struct {
		name         string
		icon         string
		expectedType string
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "empty icon",
			icon:         "",
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "stock icon",
			icon:         "warning",
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "valid file path",
			icon:         testFile,
			expectedType: "[]byte",
			expectError:  false,
		},
		{
			name:        "non-existent file",
			icon:        nonExistentFile,
			expectError: true,
			errorMsg:    "cannot read icon file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processIcon(tt.icon)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				
				if tt.expectedType == "string" {
					assert.IsType(t, "", result)
				} else if tt.expectedType == "[]byte" {
					assert.IsType(t, []byte{}, result)
					if tt.icon != "" {
						assert.Equal(t, testContent, result)
					}
				}
			}
		})
	}
}

func TestNotifyWithFileIcon(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test icon file
	iconFile := filepath.Join(tempDir, "test.png")
	iconContent := []byte("test icon content")
	err := os.WriteFile(iconFile, iconContent, 0644)
	require.NoError(t, err)
	
	mockBeeper := setupMockBeeper(t)
	mockBeeper.On("Notify", "Test", "Message", iconContent).Return(nil).Once()
	
	err = Notify("Test", "Message", iconFile, "")
	assert.NoError(t, err)
	
	mockBeeper.AssertExpectations(t)
}

func TestNotifyWithInvalidIcon(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.png")
	
	mockBeeper := setupMockBeeper(t)
	// Mock should not be called since icon processing fails
	
	err := Notify("Test", "Message", nonExistentFile, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to process icon")
	
	mockBeeper.AssertExpectations(t)
}

func TestGetSupportedFormats(t *testing.T) {
	expected := []string{".png", ".jpg", ".jpeg", ".ico", ".bmp"}
	actual := GetSupportedFormats()
	assert.Equal(t, expected, actual)
}

func TestIsValidIconFormat(t *testing.T) {
	tests := []struct {
		ext   string
		valid bool
	}{
		{".png", true},
		{".jpg", true},
		{".jpeg", true},
		{".ico", true},
		{".bmp", true},
		{".gif", false},
		{".txt", false},
		{".PNG", false}, // case sensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := IsValidIconFormat(tt.ext)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestSetBeeper(t *testing.T) {
	originalBeeper := defaultBeeper
	
	mockBeeper := new(MockBeeper)
	SetBeeper(mockBeeper)
	
	assert.Equal(t, mockBeeper, defaultBeeper)
	
	// Reset
	SetBeeper(originalBeeper)
	assert.Equal(t, originalBeeper, defaultBeeper)
}

func TestDefaultBeeper(t *testing.T) {
	beeper := NewDefaultBeeper()
	assert.NotNil(t, beeper)
	assert.IsType(t, &DefaultBeeper{}, beeper)
}