# Hanacaraka API - Onion Architecture

This document describes the onion architecture implementation for the Hanacaraka API.

## Architecture Overview

The application follows the **Onion Architecture** pattern, which provides clear separation of concerns and dependency inversion. The architecture consists of the following layers:

```
┌─────────────────────────────────────┐
│           Interfaces Layer          │  ← External concerns (HTTP handlers, middleware)
├─────────────────────────────────────┤
│         Infrastructure Layer        │  ← Data persistence, external services
├─────────────────────────────────────┤
│          Application Layer          │  ← Business logic, use cases
├─────────────────────────────────────┤
│            Domain Layer             │  ← Core entities, business rules
└─────────────────────────────────────┘
```

## Directory Structure

```
hanacaraka/
├── domain/                     # Core business layer
│   ├── entities/               # Business entities
│   │   └── user.go            # User entity with business methods
│   └── repositories/          # Repository interfaces
│       └── user_repository.go # UserRepository interface
├── application/               # Application business logic
│   └── services/             # Application services
│       └── user_service.go   # User business logic
├── infrastructure/           # External concerns
│   └── persistence/         # Data storage implementations
│       └── memory_user_repository.go # In-memory repository
├── interfaces/              # External interfaces
│   └── http/               # HTTP interface layer
│       ├── handlers/       # HTTP request handlers
│       │   └── user_handler.go # User HTTP handlers
│       └── middleware/     # HTTP middleware
│           └── logging.go  # Logging middleware
└── main.go                # Application entry point
```

## Layer Responsibilities

### Domain Layer (`domain/`)

The **innermost layer** contains:

- **Entities**: Core business objects with their properties and business methods
- **Repository Interfaces**: Contracts for data access without implementation details

**Key Principles:**

- No dependencies on external layers
- Contains pure business logic
- Defines contracts (interfaces) for external dependencies

**Files:**

- `domain/entities/user.go` - User entity with validation and business methods
- `domain/repositories/user_repository.go` - Repository interface contract

### Application Layer (`application/`)

Contains **application-specific business logic**:

- **Services**: Orchestrate business operations and enforce business rules
- **Use Cases**: Application-specific workflows

**Key Principles:**

- Depends only on the Domain layer
- Implements business workflows
- Coordinates between different domain objects

**Files:**

- `application/services/user_service.go` - User business logic and validation

### Infrastructure Layer (`infrastructure/`)

Contains **implementation details** for external concerns:

- **Database implementations**
- **External service integrations**
- **File system operations**

**Key Principles:**

- Implements interfaces defined in Domain layer
- Contains framework-specific code
- Can be easily swapped without affecting business logic

**Files:**

- `infrastructure/persistence/memory_user_repository.go` - In-memory implementation of UserRepository

### Interfaces Layer (`interfaces/`)

The **outermost layer** handles external communication:

- **HTTP handlers**
- **CLI interfaces**
- **Message queue consumers**

**Key Principles:**

- Translates external requests to application layer calls
- Handles request/response formatting
- Contains presentation logic

**Files:**

- `interfaces/http/handlers/user_handler.go` - HTTP request handlers for user operations
- `interfaces/http/middleware/logging.go` - HTTP logging middleware

## Dependency Flow

Dependencies flow **inward** only:

```
Interfaces → Infrastructure → Application → Domain
```

- **Domain** has no dependencies
- **Application** depends only on Domain
- **Infrastructure** depends on Domain (implements its interfaces)
- **Interfaces** depends on Application and Infrastructure

## Benefits

1. **Testability**: Each layer can be tested in isolation
2. **Maintainability**: Clear separation makes code easier to understand and modify
3. **Flexibility**: External dependencies can be swapped without affecting core business logic
4. **Independence**: Business logic is not tied to frameworks, databases, or external services

## API Endpoints

The application exposes the following REST endpoints:

- `GET /` - Welcome message
- `GET /users` - Get all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user
- `GET /api/v1/health` - Health check

## Running the Application

1. **Build the application:**

   ```bash
   go build
   ```

2. **Run the application:**

   ```bash
   ./com.hanacaraka
   ```

3. **Test the API:**
   ```bash
   ./test_api.sh
   ```

## Example Usage

```bash
# Get all users
curl http://localhost:8080/users

# Create a new user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Get specific user
curl http://localhost:8080/users/1

# Update user
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Smith","email":"johnsmith@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/users/1
```

## Future Enhancements

- Add database persistence (PostgreSQL, MySQL)
- Implement authentication and authorization
- Add validation middleware
- Implement proper error handling and logging
- Add unit and integration tests
- Add API documentation with Swagger
- Implement CORS middleware
- Add rate limiting
