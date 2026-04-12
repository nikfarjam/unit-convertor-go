package api

import (
	"encoding/json"
	"fmt"
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

var (
	version string = ""
	mu      sync.RWMutex
)

// WriteJSONError sends a JSON-formatted error response with correct Content-Type
func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error":"%s"}`+"\n", message)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	// Security Headers: Content-Type and nosniff to prevent MIME-sniffing vulnerabilities
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	defer r.Body.Close()
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)

	resp := VersionResponse{
		Version: loadVersion(),
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		WriteJSONError(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// versionRegex validates semantic versioning format (e.g., v1.0.0, 1.2.3-beta)
var versionRegex = regexp.MustCompile(`^v?\d+(\.\d+)*(-[\w\.-]+)?$`)

func loadVersion() string {
	mu.RLock()
	if version != "" {
		defer mu.RUnlock()
		return version
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if version != "" {
		return version
	}

	versionPath := os.Getenv("UC_VERSION_PATH")
	if versionPath == "" {
		versionPath = "version"
	}

	versionValue, err := os.ReadFile(versionPath)
	if err != nil {
		slog.Error("Error: not able to read version file", "error", err)
		version = "Unknown"
		return version
	}

	trimmedVersion := strings.TrimSpace(string(versionValue))
	// Input Validation: Ensure version string follows expected format to prevent injection or malformed data
	if !versionRegex.MatchString(trimmedVersion) {
		slog.Error("Error: invalid version format", "version", trimmedVersion)
		version = "Unknown"
		return version
	}

	version = trimmedVersion
	return version
}
