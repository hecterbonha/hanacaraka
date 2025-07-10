package routes

import (
	"net/http"

	"com.hanacaraka/application/services"
	"com.hanacaraka/interfaces/http/handlers"
	"github.com/gorilla/mux"
)

// RouteRegistry manages all route groups and their registration
type RouteRegistry struct {
	router *mux.Router
}

// NewRouteRegistry creates a new route registry
func NewRouteRegistry(router *mux.Router) *RouteRegistry {
	return &RouteRegistry{
		router: router,
	}
}

// RegisterAllRoutes registers all application routes
func (rr *RouteRegistry) RegisterAllRoutes(userService services.UserServiceInterface) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	// productHandler := handlers.NewProductHandler(productService) // Uncomment when implementing
	// orderHandler := handlers.NewOrderHandler(orderService)       // Uncomment when implementing
	// authHandler := handlers.NewAuthHandler(authService)          // Uncomment when implementing

	// Register route groups
	rr.registerBasicRoutes()
	rr.registerUserRoutes(userHandler)
	rr.registerAPIRoutes()
	rr.registerStaticRoutes()

	// Uncomment these as you implement new features
	// rr.registerProductRoutes(productHandler)
	// rr.registerOrderRoutes(orderHandler)
	// rr.registerAuthRoutes(authHandler)
}

// registerBasicRoutes registers basic application routes
func (rr *RouteRegistry) registerBasicRoutes() {
	rr.router.HandleFunc("/", homeHandler).Methods("GET")
	rr.router.HandleFunc("/ping", pingHandler).Methods("GET")
}

// registerUserRoutes registers user-related routes
func (rr *RouteRegistry) registerUserRoutes(userHandler *handlers.UserHandler) {
	SetupUserRoutes(rr.router, userHandler)
}

// registerAPIRoutes registers versioned API routes
func (rr *RouteRegistry) registerAPIRoutes() {
	// Setup API v1 routes
	apiV1 := SetupAPIRoutes(rr.router)

	// You can add API-specific route groups here
	// SetupUserAPIRoutes(apiV1, userHandler)
	// SetupProductAPIRoutes(apiV1, productHandler)

	// Setup API v2 routes (for future use)
	// apiV2 := SetupAPIV2Routes(rr.router)
}

// registerProductRoutes registers product-related routes (example for future implementation)
func (rr *RouteRegistry) registerProductRoutes(productHandler *handlers.ProductHandler) {
	// Uncomment when ProductHandler is implemented
	// SetupProductRoutes(rr.router, productHandler)
}

// registerOrderRoutes registers order-related routes (example for future implementation)
func (rr *RouteRegistry) registerOrderRoutes(orderHandler *handlers.OrderHandler) {
	// Uncomment when OrderHandler is implemented
	// router.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET")
	// router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	// router.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	// router.HandleFunc("/orders/{id}/status", orderHandler.UpdateOrderStatus).Methods("PUT")
	// router.HandleFunc("/orders/{id}/cancel", orderHandler.CancelOrder).Methods("PUT")
}

// registerAuthRoutes registers authentication-related routes (example for future implementation)
func (rr *RouteRegistry) registerAuthRoutes(authHandler *handlers.AuthHandler) {
	// Uncomment when AuthHandler is implemented
	// rr.router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	// rr.router.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	// rr.router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	// rr.router.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	// rr.router.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword).Methods("POST")
	// rr.router.HandleFunc("/auth/reset-password", authHandler.ResetPassword).Methods("POST")
}

// registerStaticRoutes registers static file serving routes
func (rr *RouteRegistry) registerStaticRoutes() {
	rr.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

// RouteGroup represents a logical grouping of routes
type RouteGroup struct {
	Name        string
	Prefix      string
	Middleware  []mux.MiddlewareFunc
	Routes      []Route
}

// Route represents a single route definition
type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
	Name    string
}

// GetRegisteredRouteGroups returns information about all registered route groups
func (rr *RouteRegistry) GetRegisteredRouteGroups() []RouteGroup {
	return []RouteGroup{
		{
			Name:   "Basic Routes",
			Prefix: "",
			Routes: []Route{
				{Path: "/", Method: "GET", Name: "Home"},
				{Path: "/ping", Method: "GET", Name: "Ping"},
			},
		},
		{
			Name:   "User Routes",
			Prefix: "/users",
			Routes: []Route{
				{Path: "/users", Method: "GET", Name: "GetUsers"},
				{Path: "/users", Method: "POST", Name: "CreateUser"},
				{Path: "/users/{id}", Method: "GET", Name: "GetUser"},
				{Path: "/users/{id}", Method: "PUT", Name: "UpdateUser"},
				{Path: "/users/{id}", Method: "DELETE", Name: "DeleteUser"},
			},
		},
		{
			Name:   "API v1 Routes",
			Prefix: "/api/v1",
			Routes: []Route{
				{Path: "/api/v1/health", Method: "GET", Name: "Health"},
				{Path: "/api/v1/status", Method: "GET", Name: "Status"},
				{Path: "/api/v1/version", Method: "GET", Name: "Version"},
			},
		},
		// Add more route groups as they are implemented
	}
}
