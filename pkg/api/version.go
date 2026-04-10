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

var versionRegex = regexp.MustCompile(`^v?\d+(\.\d+)*(-[\w\.-]+)?$`)

var cacheVersion string = ""

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)
	if cacheVersion == "" {
		cacheVersion = loadVersion()
	}
	resp := VersionResponse{
		Version: cacheVersion,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		http.Error(w, "not able to process request", http.StatusInternalServerError)
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
