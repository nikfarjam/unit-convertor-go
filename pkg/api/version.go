package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

type VersionResponse struct {
	Version string `json:"version"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var versionRegex = regexp.MustCompile(`^v?\d+(\.\d+)*(-[\w\.-]+)?$`)

var (
	cacheVersion string = ""
	mu           sync.RWMutex
)

// WriteJSONError sends a JSON-formatted error response with standard security headers.
func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		slog.Error("Error: not able to encode error response", "error", err)
	}
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)

	mu.RLock()
	v := cacheVersion
	mu.RUnlock()

	if v == "" {
		mu.Lock()
		if cacheVersion == "" {
			cacheVersion = loadVersion()
		}
		v = cacheVersion
		mu.Unlock()
	}

	resp := VersionResponse{
		Version: v,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		WriteJSONError(w, "not able to process request", http.StatusInternalServerError)
		return
	}
}

func loadVersion() string {
	versionPath := os.Getenv("UC_VERSION_PATH")
	if versionPath == "" {
		versionPath = "version"
	}
	versionValue, err := os.ReadFile(versionPath)
	if err != nil {
		slog.Error("Error: not able to read version file", "error", err)
		return "Unknown"
	}
	version := strings.TrimSpace(string(versionValue))
	if !versionRegex.MatchString(version) {
		slog.Error("Error: version format is invalid", "version", version)
		return "Unknown"
	}
	return version
}
