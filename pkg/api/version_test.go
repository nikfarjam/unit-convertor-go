package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testVersion = "1.2.3"
var test_api_url = "/version"
var default_version = "Unknown"

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

	req := httptest.NewRequest(http.MethodGet, test_api_url, nil)
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

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", w.Header().Get("Content-Type"))
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("expected X-Content-Type-Options nosniff, got %q", w.Header().Get("X-Content-Type-Options"))
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
	req = httptest.NewRequest(http.MethodGet, test_api_url, nil)
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

	var expected = default_version
	if result != expected {
		t.Errorf("%s: expected %q, got %q",
			"should return Unknown when version file not found in current dir", expected, result)
	}
}

func TestLoadVersion(t *testing.T) {
	cacheVersion = ""
	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result := loadVersion()

	if result != testVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read version from UC_VERSION_PATH", testVersion, result)
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

	if secondCall != default_version {
		t.Errorf("second call should return Unknown due to invalid path, got %q", secondCall)
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionWithWhitespace(t *testing.T) {
	cacheVersion = ""
	versionWithWhitespace := "  " + testVersion + "\n\n"
	versionFilePath := createTempVersionFile(versionWithWhitespace)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result := loadVersion()

	if result != testVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read trimmed version from UC_VERSION_PATH", testVersion, result)
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionLoadDefault(t *testing.T) {
	cacheVersion = ""
	t.Setenv("UC_VERSION_PATH", "")

	result := loadVersion()

	if result != default_version {
		t.Errorf("%s: expected %q, got %q",
			"should read trimmed version from UC_VERSION_PATH", default_version, result)
	}

}

func TestLoadVersionInvalidFormat(t *testing.T) {
	cacheVersion = ""
	invalidVersion := "harmful.sh"
	versionFilePath := createTempVersionFile(invalidVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result := loadVersion()

	if result != default_version {
		t.Errorf("%s: expected %q, got %q",
			"should read trimmed version from UC_VERSION_PATH", default_version, result)
	}
	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionInvalid(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "XSS script tag",
			content:  "<script>alert(1)</script>",
			expected: "Unknown",
		},
		{
			name:     "Invalid characters",
			content:  "1.2.3; drop table users",
			expected: "Unknown",
		},
		{
			name:     "Valid version with spaces",
			content:  "  v1.2.3  \n",
			expected: "v1.2.3",
		},
		{
			name:     "Valid semantic version with hyphen",
			content:  "1.0.0-beta.1",
			expected: "1.0.0-beta.1",
		},
		{
			name:     "Path traversal",
			content:  "cat /etc/passwd",
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheVersion = ""
			versionFilePath := createTempVersionFile(tt.content)
			t.Setenv("UC_VERSION_PATH", versionFilePath)

			result := loadVersion()

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}

			deleteTempVersionFile(versionFilePath)
		})
	}
}
