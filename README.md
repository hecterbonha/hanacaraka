# Hanacaraka HTTP API with Gorilla Mux

A simple REST API built with Go and Gorilla Mux router demonstrating HTTP routing patterns and best practices.

## Features

- RESTful API endpoints for user management
- HTTP routing with path parameters and method-specific handlers
- JSON request/response handling
- Middleware for request logging
- API versioning example
- Static file serving
- Input validation and error handling

## Prerequisites

- Go 1.24.5 or later
- curl (for testing)
- jq (optional, for pretty JSON output)

## Installation

1. Clone or navigate to the project directory
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Server

Start the server:

```bash
go run main.go
```

The server will start on port 8080 and display available endpoints.

## API Endpoints

### Core Endpoints

| Method | Endpoint      | Description       |
| ------ | ------------- | ----------------- |
| GET    | `/`           | Home page         |
| GET    | `/users`      | Get all users     |
| POST   | `/users`      | Create a new user |
| GET    | `/users/{id}` | Get user by ID    |
| PUT    | `/users/{id}` | Update user by ID |
| DELETE | `/users/{id}` | Delete user by ID |

### Additional Endpoints

| Method | Endpoint         | Description         |
| ------ | ---------------- | ------------------- |
| GET    | `/api/v1/health` | Health check        |
| GET    | `/static/*`      | Static file serving |

## Usage Examples

### Get all users

```bash
curl http://localhost:8080/users
```

### Get user by ID

```bash
curl http://localhost:8080/users/1
```

### Create a new user

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","email":"alice@example.com"}'
```

### Update a user

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'
```

### Delete a user

```bash
curl -X DELETE http://localhost:8080/users/1
```

### Health check

```bash
curl http://localhost:8080/api/v1/health
```

## Testing

Run the comprehensive test script:

```bash
./test_api.sh
```

This script tests all endpoints and demonstrates various HTTP methods and status codes.

## Project Structure

```
hanacaraka/
├── main.go          # Main application with HTTP routes
├── go.mod           # Go module file
├── test_api.sh      # API testing script
└── README.md        # This file
```

## Gorilla Mux Features Demonstrated

### 1. Basic Routing

```go
r.HandleFunc("/users", getUsersHandler).Methods("GET")
```

### 2. Path Parameters with Regex

```go
r.HandleFunc("/users/{id:[0-9]+}", getUserHandler).Methods("GET")
```

### 3. HTTP Method Restrictions

```go
r.HandleFunc("/users", createUserHandler).Methods("POST")
```

### 4. Subrouters for API Versioning

```go
api := r.PathPrefix("/api/v1").Subrouter()
api.HandleFunc("/health", healthHandler).Methods("GET")
```

### 5. Static File Serving

```go
r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
```

### 6. Middleware

```go
r.Use(loggingMiddleware)
```

### 7. Extracting Path Variables

```go
vars := mux.Vars(r)
id := vars["id"]
```

## Data Model

The API uses a simple User model:

```go
type User struct {
    ID    string `json:"id"`    // UUID string
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

The `ID` field uses UUID (Universally Unique Identifier) strings for better scalability and uniqueness guarantees. UUIDs are automatically generated when creating new users.

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET requests
- `201 Created` - Successful POST requests
- `204 No Content` - Successful DELETE requests
- `400 Bad Request` - Invalid input data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server errors

## Next Steps

To extend this API, consider adding:

- Database integration (PostgreSQL, MongoDB, etc.)
- Authentication and authorization
- Request rate limiting
- CORS middleware
- Input validation with a library like `go-playground/validator`
- Swagger/OpenAPI documentation
- Docker containerization
- Unit and integration tests
- Configuration management
- Graceful shutdown handling

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) v1.8.0 - HTTP router and URL matcher
