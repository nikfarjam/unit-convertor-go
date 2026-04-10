package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testVersion = "1.2.3"

func createTempVersionFile(content string) string {
	cacheVersion = ""

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
	cacheVersion = ""
	versionFilePath := createTempVersionFile(testVersion)

	t.Setenv("UC_VERSION_PATH", versionFilePath)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp VersionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Version != testVersion {
		t.Fatalf("expected version %q, got %q", testVersion, resp.Version)
	}

	deleteTempVersionFile(versionFilePath)
}

func TestVersionHandlerCache(t *testing.T) {
	cacheVersion = ""

	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	// First call: loads from file
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	VersionHandler(w, req)

	var resp VersionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	firstValue := resp.Version

	// Change env var to invalid path
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	// Second call: should return cached value, not try to read from new path
	req = httptest.NewRequest(http.MethodGet, "/version", nil)
	w = httptest.NewRecorder()

	VersionHandler(w, req)

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	secondValue := resp.Version

	if firstValue != secondValue {
		t.Errorf("cached value should be returned on second call: first %q, second %q", firstValue, secondValue)
	}
}

func TestLoadVersionNotFound(t *testing.T) {
	cacheVersion = ""
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	result := loadVersion()

	var expected = "Unknown"
	if string(result) != expected {
		t.Errorf("%s: expected %q, got %q",
			"should return Unknown when version file not found in current dir", expected, string(result))
	}
}

func TestLoadVersion(t *testing.T) {
	cacheVersion = ""
	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result := loadVersion()

	if string(result) != testVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read version from UC_VERSION_PATH", testVersion, string(result))
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionDoNotCache(t *testing.T) {
	cacheVersion = ""

	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	// First call: loads from file
	firstCall := loadVersion()

	// Change env var to invalid path
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	// Second call: should return cached value, not try to read from new path
	secondCall := loadVersion()

	if firstCall != testVersion {
		t.Errorf("first call should read version from file: got %q, want %q", firstCall, testVersion)
	}

	if secondCall != "Unknown" {
		t.Errorf("second call should return Unknown due to invalid path, got %q", secondCall)
	}

	deleteTempVersionFile(versionFilePath)
}
