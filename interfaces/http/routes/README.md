# Routes Organization

This directory contains the route configuration and organization for the Hanacaraka API. The routes are organized using a modular approach that separates concerns and makes it easy to add new endpoints.

## Directory Structure

```
routes/
├── README.md           # This documentation file
├── router.go          # Main router configuration
├── registry.go        # Route registry for managing all route groups
├── api_routes.go      # Versioned API routes and basic handlers
├── user_routes.go     # User-specific routes
└── product_routes.go  # Example product routes (template for new endpoints)
```

## Key Components

### Router (`router.go`)
The main router that wraps Gorilla Mux and provides the entry point for route setup. It uses the registry pattern to organize all routes.

### Registry (`registry.go`)
The route registry manages all route groups and their registration. This is where you register new route groups as you add features.

### Route Files
Each feature/domain has its own route file:
- `user_routes.go` - User CRUD operations
- `api_routes.go` - API versioning and health endpoints
- `product_routes.go` - Example for future product features

## Adding New Endpoints

### Option 1: Add to Existing Route Group

If your new endpoint belongs to an existing domain (e.g., users), add it to the appropriate route file:

```go
// In user_routes.go
func SetupUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
    // Existing routes...

    // Add your new endpoint
    router.HandleFunc("/users/{id}/profile", userHandler.GetUserProfile).Methods("GET")
    router.HandleFunc("/users/search", userHandler.SearchUsers).Methods("GET")
}
```

### Option 2: Create New Route Group

For new domains/features, create a new route file:

1. **Create the route file** (e.g., `order_routes.go`):
```go
package routes

import (
    "com.hanacaraka/interfaces/http/handlers"
    "github.com/gorilla/mux"
)

func SetupOrderRoutes(router *mux.Router, orderHandler *handlers.OrderHandler) {
    router.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET")
    router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
    router.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
    router.HandleFunc("/orders/{id}/status", orderHandler.UpdateOrderStatus).Methods("PUT")
}
```

2. **Register in the registry** (`registry.go`):
```go
func (rr *RouteRegistry) RegisterAllRoutes(userService services.UserServiceInterface, orderService services.OrderServiceInterface) {
    // Existing registrations...

    // Add new route group
    orderHandler := handlers.NewOrderHandler(orderService)
    rr.registerOrderRoutes(orderHandler)
}

func (rr *RouteRegistry) registerOrderRoutes(orderHandler *handlers.OrderHandler) {
    SetupOrderRoutes(rr.router, orderHandler)
}
```

3. **Update main.go** to pass the new service:
```go
// In main.go
orderService := services.NewOrderService(orderRepo)
router.SetupRoutes(userService, orderService)
```

## Route Patterns

### REST API Routes
Follow RESTful conventions:
```
GET    /users          # List all users
POST   /users          # Create new user
GET    /users/{id}     # Get specific user
PUT    /users/{id}     # Update specific user
DELETE /users/{id}     # Delete specific user
```

### Nested Resources
For related resources:
```
GET    /users/{id}/orders     # Get user's orders
POST   /users/{id}/orders     # Create order for user
GET    /orders/{id}/items     # Get order items
```

### API Versioning
Use path-based versioning:
```
/api/v1/users          # Version 1 API
/api/v2/users          # Version 2 API
```

### Search and Filtering
Use query parameters:
```
GET /users?name=john&status=active
GET /products?category=electronics&price_min=100
```

## Middleware

Global middleware is applied in `router.go`:
```go
r.mux.Use(middleware.LoggingMiddleware)
r.mux.Use(middleware.CORSMiddleware)
```

Route-specific middleware can be applied to subrouters:
```go
api := r.mux.PathPrefix("/api/v1").Subrouter()
api.Use(middleware.AuthMiddleware)
```

## Best Practices

1. **Group Related Routes**: Keep related endpoints in the same file
2. **Use Consistent Naming**: Follow REST conventions
3. **Version Your APIs**: Use `/api/v1/` prefix for versioned endpoints
4. **Separate Concerns**: Keep route definitions separate from business logic
5. **Document Routes**: Add comments for complex routes
6. **Use HTTP Methods Correctly**: GET for reading, POST for creating, PUT for updating, DELETE for removing

## Example Usage

Here's how routes are organized in practice:

```go
// Basic routes
GET  /                  # Welcome page
GET  /ping             # Health check

// User routes
GET    /users          # List users
POST   /users          # Create user
GET    /users/{id}     # Get user
PUT    /users/{id}     # Update user
DELETE /users/{id}     # Delete user

// API routes
GET /api/v1/health     # API health check
GET /api/v1/status     # API status
GET /api/v1/version    # API version info

// Static files
GET /static/*          # Serve static files
```

## Testing Routes

You can test routes using the provided `test_api.sh` script or curl:

```bash
# Test basic routes
curl http://localhost:8080/
curl http://localhost:8080/ping

# Test API routes
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/users

# Test user CRUD
curl -X POST http://localhost:8080/users -d '{"name":"John","email":"john@example.com"}'
curl http://localhost:8080/users/1
```

## Future Enhancements

Consider these patterns as your application grows:

1. **Route Groups with Middleware**: Apply specific middleware to route groups
2. **Route Parameters Validation**: Validate route parameters before hitting handlers
3. **Rate Limiting**: Apply rate limiting to specific route groups
4. **Authentication**: Protect routes with authentication middleware
5. **Documentation**: Auto-generate API documentation from route definitions
