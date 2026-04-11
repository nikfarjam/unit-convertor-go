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

var versionRegex = regexp.MustCompile(`^v?\d+(\.\d+)*(-[\w\.-]+)?$`)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
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
	} else {
		v := strings.TrimSpace(string(versionValue))
		if versionRegex.MatchString(v) {
			version = v
		} else {
			slog.Error("Error: invalid version format", "version", v)
			version = "Unknown"
		}
	}
	return version
}
