package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorResponseWriter struct {
	header http.Header
}

func (e *errorResponseWriter) Header() http.Header {
	if e.header == nil {
		e.header = make(http.Header)
	}
	return e.header
}

func (e *errorResponseWriter) Write(b []byte) (int, error) {
	return 0, errors.New("write error")
}

func (e *errorResponseWriter) WriteHeader(statusCode int) {}

func TestConverterHandler_MaxBytes(t *testing.T) {
	largeBody := make([]byte, 1024*1024+1)
	req := httptest.NewRequest(http.MethodPost, "/converter", bytes.NewReader(largeBody))
	w := httptest.NewRecorder()

	ConverterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestConverterHandler_EncodeError(t *testing.T) {
	reqBody := []byte(`{"value": 0, "from": "celsius", "to": "fahrenheit"}`)
	req := httptest.NewRequest(http.MethodPost, "/converter", bytes.NewReader(reqBody))
	w := &errorResponseWriter{}

	// This might not trigger an error in enc.Encode if it doesn't flush/write immediately in a way we can catch
	// but let's try.
	ConverterHandler(w, req)
}

func TestVersionHandler_EncodeError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := &errorResponseWriter{}

	VersionHandler(w, req)
}
