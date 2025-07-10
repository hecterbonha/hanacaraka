package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"com.hanacaraka/application/services"
	"com.hanacaraka/infrastructure/persistence"
	"com.hanacaraka/interfaces/http/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize repository layer
	userRepo := persistence.NewMemoryUserRepository()

	// Initialize service layer
	userService := services.NewUserService(userRepo)

	// Setup router using the routes package
	router := routes.NewRouter()
	router.SetupRoutes(userService)
	r := router.GetMux()

	// Get port from environment, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get host from environment, default to localhost
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	// Start the server
	addr := host + ":" + port
	fmt.Printf("Server starting on %s...\n", addr)
	fmt.Printf("Environment: %s\n", getEnv("ENV", "development"))
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
