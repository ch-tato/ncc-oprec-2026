package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType == expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v\nRaw Body: %s", err, rr.Body.String())
	}

	expectedKeys := []string{"name", "nrp", "status", "message", "timestamp", "uptime"}
	for _, key := range expectedKeys {
		if _, ok := responseBody[key]; !ok {
			t.Errorf("JSON response is missing expected key: '%s'", key)
		}
	}

}

func TestHealthHandler_WrongMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/health", nil)
	if err == nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code for POST request: got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}
