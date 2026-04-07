package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		envValue string
		expected slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"debug", slog.LevelDebug},
		{"DeBug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"info", slog.LevelInfo},
		{"WARN", slog.LevelWarn},
		{"warn", slog.LevelWarn},
		{"ERROR", slog.LevelError},
		{"error", slog.LevelError},
		{"", slog.LevelInfo},
		{"INVALID", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.envValue, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tt.envValue)
			result := getLogLevel()
			if result != tt.expected {
				t.Errorf("for LOG_LEVEL=%q: expected %v, got %v", tt.envValue, tt.expected, result)
			}
		})
	}
}

func TestGetLogWriter(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		logOutput    string
		logFilePath  string
		expectStdout bool
		expectedPath string
	}{
		{
			name:         "default to stdout",
			logOutput:    "",
			logFilePath:  "",
			expectStdout: true,
		},
		{
			name:         "stdout when not FILE",
			logOutput:    "CONSOLE",
			logFilePath:  "",
			expectStdout: true,
		},
		{
			name:         "file with default name",
			logOutput:    "FILE",
			logFilePath:  "",
			expectStdout: false,
			expectedPath: "", // Will create app.log in current dir, don't check
		},
		{
			name:         "file with specific path",
			logOutput:    "FILE",
			logFilePath:  filepath.Join(tempDir, "test.log"),
			expectStdout: false,
			expectedPath: filepath.Join(tempDir, "test.log"),
		},
		{
			name:         "file in directory",
			logOutput:    "FILE",
			logFilePath:  filepath.Join(tempDir, "logs/"),
			expectStdout: false,
			expectedPath: filepath.Join(tempDir, "logs/app.log"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing env
			os.Unsetenv("LOG_OUTPUT")
			os.Unsetenv("LOG_FILE_PATH")

			if tt.logOutput != "" {
				t.Setenv("LOG_OUTPUT", tt.logOutput)
			}
			if tt.logFilePath != "" {
				t.Setenv("LOG_FILE_PATH", tt.logFilePath)
			}

			writer := getLogWriter()
			if tt.expectStdout {
				if writer != os.Stdout {
					t.Errorf("expected os.Stdout, got %v", writer)
				}
			} else {
				if writer == os.Stdout {
					t.Errorf("expected file, got os.Stdout")
				}
				writer.Close() // Close the file
				if tt.expectedPath != "" {
					if _, err := os.Stat(tt.expectedPath); os.IsNotExist(err) {
						t.Errorf("expected file %s to exist", tt.expectedPath)
					}
				}
			}
		})
	}
}
