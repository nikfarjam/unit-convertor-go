package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testVersion = "1.2.3"
var testCommit = "a3f8b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9"
var test_api_url = "/version"
var default_version = "Unknown"
var default_commit = "develop"

func createTempVersionFile(content string) string {
	cacheVersion = ""
	cacheCommit = ""

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
	if resp.Commit != default_commit {
		t.Fatalf("expected commit %q (single-line file), got %q", default_commit, resp.Commit)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", w.Header().Get("Content-Type"))
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("expected X-Content-Type-Options nosniff, got %q", w.Header().Get("X-Content-Type-Options"))
	}

	deleteTempVersionFile(versionFilePath)
}

func TestVersionHandlerWithCommit(t *testing.T) {
	cacheVersion = ""
	content := testVersion + "\n" + testCommit
	versionFilePath := createTempVersionFile(content)
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
	if resp.Commit != testCommit {
		t.Fatalf("expected commit %q, got %q", testCommit, resp.Commit)
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

	result, commit := loadVersion()

	if result != default_version {
		t.Errorf("%s: expected %q, got %q",
			"should return Unknown when version file not found in current dir", default_version, result)
	}
	if commit != default_commit {
		t.Errorf("expected default commit %q, got %q", default_commit, commit)
	}
}

func TestLoadVersion(t *testing.T) {
	cacheVersion = ""
	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	result, commit := loadVersion()

	if result != testVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read version from UC_VERSION_PATH", testVersion, result)
	}
	if commit != default_commit {
		t.Errorf("expected default commit %q for single-line file, got %q", default_commit, commit)
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionDoNotCache(t *testing.T) {
	cacheVersion = ""

	versionFilePath := createTempVersionFile(testVersion)
	t.Setenv("UC_VERSION_PATH", versionFilePath)

	// First call: loads from file
	firstCall, _ := loadVersion()

	// Change env var to invalid path
	t.Setenv("UC_VERSION_PATH", "does_not_exist_version_file")

	// Second call: should return cached value, not try to read from new path
	secondCall, _ := loadVersion()

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

	result, _ := loadVersion()

	if result != testVersion {
		t.Errorf("%s: expected %q, got %q",
			"should read trimmed version from UC_VERSION_PATH", testVersion, result)
	}

	deleteTempVersionFile(versionFilePath)
}

func TestLoadVersionLoadDefault(t *testing.T) {
	cacheVersion = ""
	t.Setenv("UC_VERSION_PATH", "")

	result, _ := loadVersion()

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

	result, _ := loadVersion()

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

			result, _ := loadVersion()

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}

			deleteTempVersionFile(versionFilePath)
		})
	}
}

func TestLoadVersionCommit(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedCommit string
	}{
		{
			name:           "Full 40-char SHA",
			content:        testVersion + "\n" + testCommit,
			expectedCommit: testCommit,
		},
		{
			name:           "Short 7-char SHA",
			content:        testVersion + "\na3f8b2c",
			expectedCommit: "a3f8b2c",
		},
		{
			name:           "No second line — defaults to develop",
			content:        testVersion,
			expectedCommit: default_commit,
		},
		{
			name:           "Invalid commit hash — defaults to develop",
			content:        testVersion + "\nnot-a-hash!",
			expectedCommit: default_commit,
		},
		{
			name:           "Commit hash with surrounding whitespace",
			content:        testVersion + "\n  " + testCommit + "  ",
			expectedCommit: testCommit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheVersion = ""
			versionFilePath := createTempVersionFile(tt.content)
			t.Setenv("UC_VERSION_PATH", versionFilePath)

			_, commit := loadVersion()

			if commit != tt.expectedCommit {
				t.Errorf("expected commit %q, got %q", tt.expectedCommit, commit)
			}

			deleteTempVersionFile(versionFilePath)
		})
	}
}
