package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var expectedVersion = "1.2.3"

func createTempVersionFile(content string) string {
	version = nil

	// Create a temp file with random name in os temp folder
	tempFile, err := os.CreateTemp("", "version_*")
	if err != nil {
		panic("failed to create temp file: " + err.Error())
	}
	defer tempFile.Close()

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		panic("failed to write to temp file: " + err.Error())
	}

	return tempFile.Name()
}

func deleteTempVersionFile(path string) {
	os.Remove(path)
}

func TestVersionHandler(t *testing.T) {
	version = nil
	versionFilePath := createTempVersionFile(expectedVersion)

	t.Setenv("UC_VERSION_PATH", versionFilePath)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if got := w.Body.String(); got != expectedVersion {
		t.Fatalf("expected body %q, got %q", expectedVersion, got)
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionNotFound(t *testing.T) {
	version = nil
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	result := loadVersion()

	var expected = "Unknown"
	if string(result) != expected {
		t.Errorf("%s: expected %q, got %q",
			"should return Unknown when version file not found in current dir", expected, string(result))
	}
}

func TestLoadVersion(t *testing.T) {
	version = nil
	versionFilePath := createTempVersionFile(expectedVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result := loadVersion()

	if string(result) != expectedVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read version from UC_VERSION_PATH", expectedVersion, string(result))
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionCaching(t *testing.T) {
	version = nil

	versionFilePath := createTempVersionFile(expectedVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

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

	deleteTempVersionFile(versionFilePath)
}
