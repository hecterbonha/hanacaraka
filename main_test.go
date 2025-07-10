package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request to home",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to Hanacaraka API!\n",
		},
		{
			name:           "POST request to home (should still work)",
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to Hanacaraka API!\n",
		},
		{
			name:           "PUT request to home (should still work)",
			method:         "PUT",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to Hanacaraka API!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(homeHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedHealth string
	}{
		{
			name:           "GET request to health",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedHealth: "healthy",
		},
		{
			name:           "POST request to health (should still work)",
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedHealth: "healthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/api/v1/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(healthHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			// Check Content-Type header
			expectedContentType := "application/json"
			if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("Expected Content-Type '%s', got '%s'", expectedContentType, contentType)
			}

			// Parse and check JSON response
			var response map[string]string
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Errorf("Error decoding JSON response: %v", err)
			}

			if status, exists := response["status"]; !exists {
				t.Error("Expected 'status' field in response")
			} else if status != tt.expectedHealth {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedHealth, status)
			}
		})
	}
}

func TestHealthHandler_JSONStructure(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	// Verify it's valid JSON
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Check that it only contains the expected field
	if len(response) != 1 {
		t.Errorf("Expected response to have exactly 1 field, got %d", len(response))
	}

	// Check that the status field exists and is a string
	if status, exists := response["status"]; !exists {
		t.Error("Expected 'status' field in response")
	} else {
		if statusStr, ok := status.(string); !ok {
			t.Errorf("Expected 'status' to be a string, got %T", status)
		} else if statusStr != "healthy" {
			t.Errorf("Expected status 'healthy', got '%s'", statusStr)
		}
	}
}

func TestHealthHandler_ResponseFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	// Check that response can be unmarshaled into expected structure
	type HealthResponse struct {
		Status string `json:"status"`
	}

	var healthResp HealthResponse
	err = json.NewDecoder(rr.Body).Decode(&healthResp)
	if err != nil {
		t.Errorf("Error unmarshaling response into HealthResponse struct: %v", err)
	}

	if healthResp.Status != "healthy" {
		t.Errorf("Expected Status 'healthy', got '%s'", healthResp.Status)
	}
}

func TestHomeHandler_ResponseFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	// Check that response ends with newline (as expected from Fprintf)
	body := rr.Body.String()
	if body[len(body)-1] != '\n' {
		t.Error("Expected response to end with newline")
	}

	// Check that response starts with expected text
	expectedPrefix := "Welcome to Hanacaraka API!"
	if !strings.HasPrefix(body, expectedPrefix) {
		t.Errorf("Expected response to start with '%s', got '%s'", expectedPrefix, body)
	}
}

func TestHandlers_NoSideEffects(t *testing.T) {
	// Test that handlers don't modify global state or have side effects

	// Test homeHandler multiple times
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(homeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Iteration %d: Expected status %d, got %d", i, http.StatusOK, status)
		}

		expectedBody := "Welcome to Hanacaraka API!\n"
		if body := rr.Body.String(); body != expectedBody {
			t.Errorf("Iteration %d: Expected body '%s', got '%s'", i, expectedBody, body)
		}
	}

	// Test healthHandler multiple times
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", "/api/v1/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Iteration %d: Expected status %d, got %d", i, http.StatusOK, status)
		}

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Errorf("Iteration %d: Error decoding JSON: %v", i, err)
		}

		if response["status"] != "healthy" {
			t.Errorf("Iteration %d: Expected status 'healthy', got '%s'", i, response["status"])
		}
	}
}
