package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestSetupServer(t *testing.T) {
	server := setupServer()

	if server.Addr != ":9090" {
		t.Errorf("expected addr :9090, got %s", server.Addr)
	}

	if server.ReadTimeout != 5*time.Second {
		t.Errorf("expected ReadTimeout 5s, got %v", server.ReadTimeout)
	}

	if server.WriteTimeout != 10*time.Second {
		t.Errorf("expected WriteTimeout 10s, got %v", server.WriteTimeout)
	}

	if server.IdleTimeout != 120*time.Second {
		t.Errorf("expected IdleTimeout 120s, got %v", server.IdleTimeout)
	}

	mux, ok := server.Handler.(*http.ServeMux)
	if !ok {
		t.Fatal("expected handler to be *http.ServeMux")
	}

	// We can't easily check registered routes in http.ServeMux without reflection or internal access
	// but we've verified the structure of setupServer.
	if mux == nil {
		t.Fatal("mux should not be nil")
	}
}

func TestInitLogger(t *testing.T) {
	// Call initLogger to cover it
	initLogger()

	// Verify that the default logger is set (JSON handler)
	// slog doesn't expose its handler easily, but we can verify it doesn't panic
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
		{
			name:         "file without .log extension",
			logOutput:    "FILE",
			logFilePath:  filepath.Join(tempDir, "myapp"),
			expectStdout: false,
			expectedPath: filepath.Join(tempDir, "myapp/app.log"),
		},
		{
			name:         "invalid path defaults to stdout",
			logOutput:    "FILE",
			logFilePath:  "/invalid/path/to/log.log",
			expectStdout: true,
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
