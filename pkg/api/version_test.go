package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var expectedVersion string

func setup() {
	version = nil
	// Load the version from the project root version file
	cwd, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}

	versionPath := filepath.Join(cwd, "..", "..", "version")
	content, err := os.ReadFile(versionPath)
	if err != nil {
		panic("failed to read version file: " + err.Error())
	}

	expectedVersion = strings.TrimSpace(string(content))
}

func TestVersionHandler(t *testing.T) {
	setup()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// The test package path is pkg/api; point to the project root version file explicitly.
	projectVersionPath := filepath.Join(cwd, "..", "..", "version")
	t.Setenv("UC_VERSION_PATH", projectVersionPath)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if got := w.Body.String(); got != expectedVersion {
		t.Fatalf("expected body %q, got %q", expectedVersion, got)
	}
}

func TestLoadVersionNotFound(t *testing.T) {
	version = nil
	t.Setenv("UC_VERSION_PATH", "")

	result := loadVersion()

	var expected = "Unknown"
	if string(result) != expected {
		t.Errorf("%s: expected %q, got %q",
			"should return Unknown when version file not found in current dir", expected, string(result))
	}
}

func TestLoadVersion(t *testing.T) {
	setup()
	cwd, _ := os.Getwd()
	versionPath := filepath.Join(cwd, "..", "..", "version")
	t.Setenv("UC_VERSION_PATH", versionPath)

	result := loadVersion()

	if string(result) != expectedVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read version from UC_VERSION_PATH", expectedVersion, string(result))
	}

}

func TestLoadVersionCaching(t *testing.T) {
	setup()

	cwd, _ := os.Getwd()
	versionPath := filepath.Join(cwd, "..", "..", "version")
	t.Setenv("UC_VERSION_PATH", versionPath)

	// First call: loads from file
	firstCall := loadVersion()

	// Change env var to invalid path
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	// Second call: should return cached value, not try to read from new path
	secondCall := loadVersion()

	if string(firstCall) != string(secondCall) {
		t.Errorf("cached value should be returned on second call: first %q, second %q", string(firstCall), string(secondCall))
	}

	if string(secondCall) != expectedVersion {
		t.Errorf("cached value should match expected version: got %q, want %q", string(secondCall), expectedVersion)
	}
}
