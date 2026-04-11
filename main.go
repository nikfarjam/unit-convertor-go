package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"strings"

	"github.com/nikfarjam/unit-convertor-go/pkg/api"
)

func setupServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /converter", api.ConverterHandler)
	mux.HandleFunc("GET /version", api.VersionHandler)

	addr := ":9090"
	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func main() {
	initLogger()

	server := setupServer()

	fmt.Printf("Server is running on http://localhost%s\n", server.Addr)
	slog.Warn("Server is running on http://localhost" + server.Addr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}

func getLogLevel() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(getLogWriter(), &slog.HandlerOptions{
		Level: getLogLevel(),
	}))
	slog.SetDefault(logger)
}

func getLogWriter() *os.File {
	logOutput := os.Getenv("LOG_OUTPUT")

	if strings.ToUpper(logOutput) != "FILE" {
		return os.Stdout
	}

	logPath := os.Getenv("LOG_FILE_PATH")
	if logPath == "" {
		logPath = "app.log"
	} else if !strings.HasSuffix(logPath, ".log") {
		if !strings.HasSuffix(logPath, "/") {
			logPath += "/"
		}
		logPath += "app.log"
	}

	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating log directory: %v, defaulting to stdout\n", err)
	} else {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v, defaulting to stdout\n", err)
		} else {
			return file
		}
	}
	return os.Stdout
}
