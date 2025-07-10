package routes

import (
	"com.hanacaraka/interfaces/http/handlers"
	"github.com/gorilla/mux"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	// User CRUD operations
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Additional user routes can be added here
	// router.HandleFunc("/users/{id}/profile", userHandler.GetUserProfile).Methods("GET")
	// router.HandleFunc("/users/{id}/avatar", userHandler.UpdateUserAvatar).Methods("PUT")
	// router.HandleFunc("/users/search", userHandler.SearchUsers).Methods("GET")
}

// SetupUserAPIRoutes configures versioned API user routes
func SetupUserAPIRoutes(apiRouter *mux.Router, userHandler *handlers.UserHandler) {
	// API v1 user routes
	apiRouter.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	apiRouter.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Additional API user endpoints
	// apiRouter.HandleFunc("/users/bulk", userHandler.BulkCreateUsers).Methods("POST")
	// apiRouter.HandleFunc("/users/export", userHandler.ExportUsers).Methods("GET")
}
