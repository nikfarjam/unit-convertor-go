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
	if version == nil {
		versionValue, err := os.ReadFile("version")
		if err != nil {
			slog.Error("Error: not able to read version file", "error", err)
			version = []byte("Unknown")
		} else {
			version = versionValue
		}
	}
	_, err := w.Write(version)
	if err != nil {
		slog.Error("Error: not able to write response", "error", err)
		http.Error(w, "not able to write response", http.StatusInternalServerError)
		return
	}
}
