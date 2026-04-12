package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nikfarjam/unit-convertor-go/pkg/converter"
)

var api_url = "/convert"

func TestConverterHandler_Success(t *testing.T) {
	reqBody := converter.ConverterRequest{
		Value: 0,
		From:  "celsius",
		To:    "fahrenheit",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp converter.ConverterResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedValue := 32.0
	if resp.Value != expectedValue {
		t.Errorf("expected value %f, got %f", expectedValue, resp.Value)
	}

	if resp.Unit != "FAHRENHEIT" {
		t.Errorf("expected unit %s, got %s", "FAHRENHEIT", resp.Unit)
	}
}

func TestConverterHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expected := `{"error":"bad request"}` + "\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

func TestConverterHandler_BodySizeLimit(t *testing.T) {
	// Create a valid JSON body larger than 1MB
	// {"value": 0, "from": "celsius", "to": "fahrenheit", "padding": "..."}
	padding := strings.Repeat("a", 1048576)
	reqBody := map[string]interface{}{
		"value":   0,
		"from":    "celsius",
		"to":      "fahrenheit",
		"padding": padding,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expected := `{"error":"bad request"}` + "\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

func TestConverterHandler_InvalidFromUnit(t *testing.T) {
	reqBody := converter.ConverterRequest{
		Value: 0,
		From:  "invalid",
		To:    "fahrenheit",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expected := `{"error":"invalid from"}` + "\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

func TestConverterHandler_InvalidToUnit(t *testing.T) {
	reqBody := converter.ConverterRequest{
		Value: 0,
		From:  "celsius",
		To:    "invalid",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expected := `{"error":"invalid to"}` + "\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

func TestConverterHandler_FahrenheitToCelsius(t *testing.T) {
	reqBody := converter.ConverterRequest{
		Value: 32,
		From:  "fahrenheit",
		To:    "celsius",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp converter.ConverterResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedValue := 0.0
	if resp.Value != expectedValue {
		t.Errorf("expected value %f, got %f", expectedValue, resp.Value)
	}

	if resp.Unit != "CELSIUS" {
		t.Errorf("expected unit %s, got %s", "CELSIUS", resp.Unit)
	}
}

func TestConverterHandler_ConversionError(t *testing.T) {
	// Identity conversion might fail if not supported (e.g., CELSIUS to CELSIUS)
	reqBody := converter.ConverterRequest{
		Value: 25,
		From:  "CELSIUS",
		To:    "CELSIUS",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expected := `{"error":"not able to process request"}` + "\n"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

// errorResponseWriter is a mock ResponseWriter that fails on Write
type errorResponseWriter struct {
	header http.Header
}

func (e *errorResponseWriter) Header() http.Header {
	if e.header == nil {
		e.header = make(http.Header)
	}
	return e.header
}
func (e *errorResponseWriter) Write(b []byte) (int, error) { return 0, fmt.Errorf("write error") }
func (e *errorResponseWriter) WriteHeader(statusCode int)  {}

func TestConverterHandler_EncodeError(t *testing.T) {
	reqBody := converter.ConverterRequest{
		Value: 25,
		From:  "CELSIUS",
		To:    "FAHRENHEIT",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := &errorResponseWriter{}

	// This will trigger the slog.Error for encoding failure but won't be able to WriteJSONError either
	// because Write also fails.
	ConverterHandler(w, req)
}
