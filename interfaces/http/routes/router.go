package routes

import (
	"com.hanacaraka/application/services"
	"com.hanacaraka/interfaces/http/middleware"

	"github.com/gorilla/mux"
)

// Router wraps the mux router and provides route setup methods
type Router struct {
	mux      *mux.Router
	registry *RouteRegistry
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	router := mux.NewRouter()
	return &Router{
		mux:      router,
		registry: NewRouteRegistry(router),
	}
}

// GetMux returns the underlying mux router
func (r *Router) GetMux() *mux.Router {
	return r.mux
}

// SetupRoutes configures all routes for the application
func (r *Router) SetupRoutes(userService services.UserServiceInterface) {
	// Apply global middleware
	r.mux.Use(middleware.LoggingMiddleware)

	// Register all routes using the registry
	r.registry.RegisterAllRoutes(userService)
}

// GetRouteRegistry returns the route registry for inspection
func (r *Router) GetRouteRegistry() *RouteRegistry {
	return r.registry
}
