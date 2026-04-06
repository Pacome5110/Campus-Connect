package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestWebhookInvalidPayload(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestWebhookValidPayload(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte(`{"test":"hi"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestMissingAPIKey(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/analytics", nil)
	router.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestWrongAPIKey(t *testing.T) {
	router := SetupRouter()
	os.Setenv("API_KEY", "real-key")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/analytics", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	router.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("Expected 401 for wrong key, got %d", w.Code)
	}
}

func TestCorrectAPIKey(t *testing.T) {
	router := SetupRouter()
	os.Setenv("API_KEY", "real-key")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/analytics", nil)
	req.Header.Set("X-API-Key", "real-key")
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected 200 for valid API key, got %d", w.Code)
	}
}

func TestNotificationWithCorrectKey(t *testing.T) {
	router := SetupRouter()
	os.Setenv("API_KEY", "real-key")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/notifications", bytes.NewBuffer([]byte(`{"message":"hello"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "real-key")
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestRateLimiter(t *testing.T) {
	router := SetupRouter()
	// Limit is 5. Doing 6 requests should give 429 on the 6th.
	// We use a unique mocked IP so other tests running in sequence don't interfere
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		req.Header.Set("X-Forwarded-For", "192.168.10.99")
		router.ServeHTTP(w, req)
	}
	
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/health", nil)
	req2.Header.Set("X-Forwarded-For", "192.168.10.99")
	router.ServeHTTP(w2, req2)

	if w2.Code != 429 {
		t.Errorf("Expected 429 Rate limit exceeded, got %d", w2.Code)
	}
}
