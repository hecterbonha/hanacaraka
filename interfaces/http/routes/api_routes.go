package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupAPIRoutes configures all versioned API routes
func SetupAPIRoutes(router *mux.Router) *mux.Router {
	// Create API v1 subrouter
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Health and status endpoints
	apiV1.HandleFunc("/health", healthHandler).Methods("GET")
	apiV1.HandleFunc("/status", statusHandler).Methods("GET")
	apiV1.HandleFunc("/version", versionHandler).Methods("GET")

	// Example: Product endpoints (when you add product functionality)
	// apiV1.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	// apiV1.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	// apiV1.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	// apiV1.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	// apiV1.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Example: Order endpoints (when you add order functionality)
	// apiV1.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET")
	// apiV1.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	// apiV1.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	// apiV1.HandleFunc("/orders/{id}/status", orderHandler.UpdateOrderStatus).Methods("PUT")

	// Example: Authentication endpoints
	// apiV1.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	// apiV1.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	// apiV1.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")

	return apiV1
}

// SetupAPIV2Routes configures API v2 routes (for future use)
func SetupAPIV2Routes(router *mux.Router) *mux.Router {
	apiV2 := router.PathPrefix("/api/v2").Subrouter()

	// Future v2 endpoints will go here
	apiV2.HandleFunc("/health", healthV2Handler).Methods("GET")

	return apiV2
}

// API route handlers
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": "v1",
		"service": "hanacaraka",
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "operational",
		"timestamp": "2024-01-01T00:00:00Z", // In real implementation, use time.Now()
		"uptime":    "24h30m",
		"version":   "v1",
	})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version":    "1.0.0",
		"api_version": "v1",
		"build":      "development",
		"commit":     "latest",
	})
}

func healthV2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": "v2",
		"service": "hanacaraka",
		"features": []string{"enhanced_logging", "metrics", "tracing"},
	})
}

// Basic route handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Hanacaraka API!\n"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
		"status":  "ok",
	})
}
