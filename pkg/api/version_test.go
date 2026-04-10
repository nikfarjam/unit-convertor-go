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
	version = ""

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
	version = ""
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

	// Verify security headers
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", w.Header().Get("Content-Type"))
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("expected X-Content-Type-Options nosniff, got %q", w.Header().Get("X-Content-Type-Options"))
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionInvalidFormat(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"Empty", ""},
		{"Not a version", "hello-world"},
		{"Path injection", "../etc/passwd"},
		{"Script tag", "<script>alert(1)</script>"},
		{"Special characters", "v1.0.0; DROP TABLE users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version = ""
			versionFilePath := createTempVersionFile(tt.content)
			t.Setenv("UC_VERSION_PATH", versionFilePath)

			result := loadVersion()
			if result != "Unknown" {
				t.Errorf("expected Unknown for content %q, got %q", tt.content, result)
			}
			deleteTempVersionFile(versionFilePath)
		})
	}
}

func TestLoadVersionNotFound(t *testing.T) {
	version = ""
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	result := loadVersion()

	var expected = "Unknown"
	if string(result) != expected {
		t.Errorf("%s: expected %q, got %q",
			"should return Unknown when version file not found in current dir", expected, string(result))
	}
}

func TestLoadVersion(t *testing.T) {
	version = ""
	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result := loadVersion()

	if string(result) != testVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read version from UC_VERSION_PATH", testVersion, string(result))
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionCaching(t *testing.T) {
	version = ""

	versionFilePath := createTempVersionFile(testVersion)
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

	if string(secondCall) != testVersion {
		t.Errorf("cached value should match expected version: got %q, want %q", string(secondCall), testVersion)
	}

	deleteTempVersionFile(versionFilePath)
}
