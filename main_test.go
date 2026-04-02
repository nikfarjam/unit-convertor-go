package main

import (
	"log/slog"
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
				t.Errorf("getLogLevel() = %v, want %v for input '%s'", result, tt.expected, tt.envValue)
			}
		})
	}
}
