package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		remoteAddr     string
		expectedLog    string
	}{
		{
			name:        "GET request",
			method:      "GET",
			url:         "/users",
			remoteAddr:  "192.168.1.1:12345",
			expectedLog: "192.168.1.1:12345 GET /users",
		},
		{
			name:        "POST request",
			method:      "POST",
			url:         "/users",
			remoteAddr:  "127.0.0.1:54321",
			expectedLog: "127.0.0.1:54321 POST /users",
		},
		{
			name:        "PUT request with ID",
			method:      "PUT",
			url:         "/users/123",
			remoteAddr:  "10.0.0.1:8080",
			expectedLog: "10.0.0.1:8080 PUT /users/123",
		},
		{
			name:        "DELETE request",
			method:      "DELETE",
			url:         "/users/456",
			remoteAddr:  "172.16.0.1:9090",
			expectedLog: "172.16.0.1:9090 DELETE /users/456",
		},
		{
			name:        "Request with query parameters",
			method:      "GET",
			url:         "/users?page=1&limit=10",
			remoteAddr:  "192.168.0.100:3000",
			expectedLog: "192.168.0.100:3000 GET /users?page=1&limit=10",
		},
		{
			name:        "Request with complex path",
			method:      "GET",
			url:         "/api/v1/users/123/profile",
			remoteAddr:  "203.0.113.1:443",
			expectedLog: "203.0.113.1:443 GET /api/v1/users/123/profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture log output
			var logOutput bytes.Buffer

			// Save original log output and set custom writer
			originalOutput := log.Writer()
			originalFlags := log.Flags()
			log.SetOutput(&logOutput)
			log.SetFlags(0) // Remove timestamp and other flags for cleaner testing

			// Restore original log settings after test
			defer func() {
				log.SetOutput(originalOutput)
				log.SetFlags(originalFlags)
			}()

			// Create a test handler that we'll wrap with the middleware
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			// Wrap the test handler with logging middleware
			wrappedHandler := LoggingMiddleware(testHandler)

			// Create a test request
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.RemoteAddr = tt.remoteAddr

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Execute the request
			wrappedHandler.ServeHTTP(rr, req)

			// Check that the underlying handler was called correctly
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, status)
			}

			expectedBody := "test response"
			if body := rr.Body.String(); body != expectedBody {
				t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
			}

			// Check the log output
			logString := strings.TrimSpace(logOutput.String())
			if !strings.Contains(logString, tt.expectedLog) {
				t.Errorf("Expected log to contain '%s', got '%s'", tt.expectedLog, logString)
			}
		})
	}
}

func TestLoggingMiddleware_CallsNextHandler(t *testing.T) {
	var handlerCalled bool
	var capturedRequest *http.Request
	var capturedWriter http.ResponseWriter

	// Create a test handler that records if it was called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedRequest = r
		capturedWriter = w
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("handler response"))
	})

	// Wrap with logging middleware
	wrappedHandler := LoggingMiddleware(testHandler)

	// Create test request
	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:8080"

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(rr, req)

	// Verify that the next handler was called
	if !handlerCalled {
		t.Error("Expected next handler to be called, but it wasn't")
	}

	// Verify that the correct request was passed
	if capturedRequest != req {
		t.Error("Expected same request to be passed to next handler")
	}

	// Verify that the correct response writer was passed
	if capturedWriter != rr {
		t.Error("Expected same response writer to be passed to next handler")
	}

	// Verify that the response from the next handler is preserved
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	expectedBody := "handler response"
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestLoggingMiddleware_HandlerPanic(t *testing.T) {
	// Create a handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap with logging middleware
	wrappedHandler := LoggingMiddleware(panicHandler)

	// Create test request
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:8080"

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request and expect panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to be propagated, but it wasn't")
		}
	}()

	wrappedHandler.ServeHTTP(rr, req)
}

func TestLoggingMiddleware_MultipleRequests(t *testing.T) {
	// Create a buffer to capture log output
	var logOutput bytes.Buffer

	// Save original log output and set custom writer
	originalOutput := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&logOutput)
	log.SetFlags(0)

	// Restore original log settings after test
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	// Create a simple test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with logging middleware
	wrappedHandler := LoggingMiddleware(testHandler)

	// Make multiple requests
	requests := []struct {
		method     string
		url        string
		remoteAddr string
	}{
		{"GET", "/users", "192.168.1.1:12345"},
		{"POST", "/users", "192.168.1.2:12346"},
		{"PUT", "/users/123", "192.168.1.3:12347"},
	}

	for _, reqData := range requests {
		req, err := http.NewRequest(reqData.method, reqData.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = reqData.remoteAddr

		rr := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rr, req)
	}

	// Check that all requests were logged
	logString := logOutput.String()

	expectedLogs := []string{
		"192.168.1.1:12345 GET /users",
		"192.168.1.2:12346 POST /users",
		"192.168.1.3:12347 PUT /users/123",
	}

	for _, expectedLog := range expectedLogs {
		if !strings.Contains(logString, expectedLog) {
			t.Errorf("Expected log to contain '%s', but it's missing from: %s", expectedLog, logString)
		}
	}
}

func TestLoggingMiddleware_EmptyRemoteAddr(t *testing.T) {
	// Create a buffer to capture log output
	var logOutput bytes.Buffer

	// Save original log output and set custom writer
	originalOutput := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&logOutput)
	log.SetFlags(0)

	// Restore original log settings after test
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with logging middleware
	wrappedHandler := LoggingMiddleware(testHandler)

	// Create request with empty RemoteAddr
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	// RemoteAddr is empty by default

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	// Check that logging still works with empty RemoteAddr
	logString := strings.TrimSpace(logOutput.String())
	if !strings.Contains(logString, "GET /test") {
		t.Errorf("Expected log to contain 'GET /test', got '%s'", logString)
	}
}

func TestLoggingMiddleware_PreservesResponseHeaders(t *testing.T) {
	// Create a handler that sets custom headers
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Custom-Header", "test-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "test"}`))
	})

	// Wrap with logging middleware
	wrappedHandler := LoggingMiddleware(testHandler)

	// Create test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1:8080"

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(rr, req)

	// Verify that headers are preserved
	if customHeader := rr.Header().Get("Custom-Header"); customHeader != "test-value" {
		t.Errorf("Expected Custom-Header 'test-value', got '%s'", customHeader)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Verify status code is preserved
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, status)
	}

	// Verify response body is preserved
	expectedBody := `{"message": "test"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}
