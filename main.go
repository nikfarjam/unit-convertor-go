package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikfarjam/unit-convertor-go/pkg/converter"
)

func main() {
	initLogger()
	http.HandleFunc("/converter", converterHandler)
	http.HandleFunc("/version", versionHandler)
	addr := ":9090"
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	slog.Warn("Server is running on http://localhost" + addr)
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
	logOutput := os.Getenv("LOG_OUTPUT")
	writer := os.Stdout
	if strings.ToUpper(logOutput) == "FILE" {
		logPath := os.Getenv("LOG_FILE_PATH")
		if logPath == "" {
			logPath = "app.log"
		}
		if !strings.HasSuffix(logPath, ".log") {
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
				writer = file
			}
		}
	}
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: getLogLevel(),
	}))
	slog.SetDefault(logger)
}

func converterHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)
	dec := json.NewDecoder(r.Body)
	req := &converter.ConverterRequest{}

	if err := dec.Decode(req); err != nil {
		slog.Error("Error: bad request", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	slog.Debug("Request received", "from", req.From, "to", req.To)

	if strings.ToUpper(req.From) != "CELSIUS" && strings.ToUpper(req.From) != "FAHRENHEIT" {
		slog.Error("Error: invalid unit", "unit", req.From)
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}

	if strings.ToUpper(req.To) != "CELSIUS" && strings.ToUpper(req.To) != "FAHRENHEIT" {
		slog.Error("Error: invalid unit", "unit", req.To)
		http.Error(w, "invalid to", http.StatusBadRequest)
		return
	}

	resp, err := converter.ConvertUnit(*req)
	if err != nil {
		slog.Error("Error: not able to process request", "error", err)
		http.Error(w, "not able to process request", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		http.Error(w, "not able to process request", http.StatusInternalServerError)
		return
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)
	version, err := os.ReadFile("version")
	if err != nil {
		slog.Error("Error: not able to read version file", "error", err)
		http.Error(w, "not able to read version file", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(version)
	if err != nil {
		slog.Error("Error: not able to write response", "error", err)
		http.Error(w, "not able to write response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
