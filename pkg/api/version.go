package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type VersionResponse struct {
	Version string `json:"version"`
}

var version string = ""

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)

	// Set security headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	resp := VersionResponse{
		Version: loadVersion(),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		http.Error(w, "not able to process request", http.StatusInternalServerError)
		return
	}
}

func loadVersion() string {
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

	v := strings.TrimSpace(string(versionValue))
	// Basic semver-like validation to prevent injection or weird data
	re := regexp.MustCompile(`^v?\d+(\.\d+)*(-[\w\.-]+)?$`)
	if !re.MatchString(v) {
		slog.Error("Error: invalid version format", "version", v)
		version = "Unknown"
	} else {
		version = v
	}

	return version
}
