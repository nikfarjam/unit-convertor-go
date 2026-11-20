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
	Commit  string `json:"commit"`
}

var versionRegex = regexp.MustCompile(`^v?\d+(\.\d+)*(-[\w\.-]+)?$`)
var commitRegex = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

const defaultCommit = "develop"

var (
	cacheVersion string = ""
	cacheCommit  string = ""
	mu           sync.RWMutex
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)

	mu.RLock()
	v := cacheVersion
	c := cacheCommit
	mu.RUnlock()

	if v == "" {
		mu.Lock()
		if cacheVersion == "" {
			cacheVersion, cacheCommit = loadVersion()
		}
		v = cacheVersion
		c = cacheCommit
		mu.Unlock()
	}

	resp := VersionResponse{
		Version: v,
		Commit:  c,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		slog.Error("Error: not able to encode response", "error", err)
		WriteJSONError(w, "not able to process request", http.StatusInternalServerError)
		return
	}
}

func loadVersion() (string, string) {
	versionPath := os.Getenv("UC_VERSION_PATH")
	if versionPath == "" {
		versionPath = "version"
	}
	versionValue, err := os.ReadFile(versionPath)
	if err != nil {
		slog.Error("Error: not able to read version file", "error", err)
		return "Unknown", defaultCommit
	}

	lines := strings.SplitN(strings.TrimRight(string(versionValue), "\r\n"), "\n", 2)

	version := strings.TrimSpace(lines[0])
	if !versionRegex.MatchString(version) {
		slog.Error("Error: version format is invalid", "version", version)
		return "Unknown", defaultCommit
	}

	commit := defaultCommit
	if len(lines) == 2 {
		candidate := strings.TrimSpace(lines[1])
		if commitRegex.MatchString(candidate) {
			commit = candidate
		} else {
			slog.Warn("Commit hash format is invalid, using default", "commit", candidate)
		}
	}

	return version, commit
}
