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

var versionRegexp = regexp.MustCompile(`^[0-9]+(\.[0-9]+)*([+-][A-Za-z0-9._-]+)?$`)

var cacheVersion string = ""

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
	if !versionRegexp.MatchString(version) {
		slog.Error("Error: version format is invalid", "version", version)
		return "Unknown"
	}
	return version
}
