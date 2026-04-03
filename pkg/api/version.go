package api

import (
	"log/slog"
	"net/http"
	"os"
)

var version []byte = nil

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	slog.Debug("Received request", "method", r.Method, "path", r.URL.Path)

	_, err := w.Write(loadVersion())
	if err != nil {
		slog.Error("Error: not able to write response", "error", err)
		http.Error(w, "not able to write response", http.StatusInternalServerError)
		return
	}
}

func loadVersion() []byte {
	if version != nil {
		return version
	}
	versionPath := os.Getenv("UC_VERSION_PATH")
	if versionPath == "" {
		versionPath = "version"
	}
	versionValue, err := os.ReadFile(versionPath)
	if err != nil {
		slog.Error("Error: not able to read version file", "error", err)
		version = []byte("Unknown")
	} else {
		version = versionValue
	}
	return version
}
