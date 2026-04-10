package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	// Verify security headers
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", w.Header().Get("Content-Type"))
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("expected X-Content-Type-Options nosniff, got %q", w.Header().Get("X-Content-Type-Options"))
	}
}

func TestConverterHandler_MaxBytes(t *testing.T) {
	// Create a body larger than 1MB
	largeBody := make([]byte, 1024*1024+1)
	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(largeBody))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestConverterHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if w.Body.String() != "bad request\n" {
		t.Errorf("expected body 'bad request', got %q", w.Body.String())
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

	if w.Body.String() != "invalid from\n" {
		t.Errorf("expected body 'invalid from', got %q", w.Body.String())
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

	if w.Body.String() != "invalid to\n" {
		t.Errorf("expected body 'invalid to', got %q", w.Body.String())
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

type errorResponseWriter struct {
	httptest.ResponseRecorder
}

func (e *errorResponseWriter) Write(b []byte) (int, error) {
	return 0, bytes.ErrTooLarge // Simulating an error
}

func TestConverterHandler_EncodeError(t *testing.T) {
	reqBody := converter.ConverterRequest{
		Value: 0,
		From:  "celsius",
		To:    "fahrenheit",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, api_url, bytes.NewReader(body))
	w := &errorResponseWriter{*httptest.NewRecorder()}

	// This is a bit tricky because json.Encoder might buffer.
	// But let's see if it hits the error path.
	ConverterHandler(w, req)
	// We don't check code because it might have already sent 200 header
}
