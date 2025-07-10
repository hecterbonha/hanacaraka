# Testing Guide for Hanacaraka API

This document provides comprehensive information about testing the Hanacaraka API project.

## Overview

The project follows a clean architecture pattern with comprehensive unit tests for each layer:

- **Domain Layer**: Entity business logic and validation
- **Application Layer**: Service layer with business rules
- **Infrastructure Layer**: Repository implementations
- **Interface Layer**: HTTP handlers and middleware

## Test Structure

### Test Files

Each Go source file has a corresponding test file with the `_test.go` suffix:

```
hanacaraka/
├── main_test.go                                    # Main package handlers
├── domain/entities/user_test.go                    # User entity tests
├── application/services/user_service_test.go       # User service tests
├── infrastructure/persistence/memory_user_repository_test.go  # Repository tests
├── interfaces/http/handlers/user_handler_test.go   # HTTP handler tests
└── interfaces/http/middleware/logging_test.go      # Middleware tests
```

### Test Coverage

Current test coverage by package:

- **Domain Entities**: 100% coverage
- **Infrastructure**: 100% coverage
- **HTTP Handlers**: 100% coverage
- **Middleware**: 100% coverage
- **Application Services**: 90.9% coverage
- **Main Package**: 15% coverage (handlers only)

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Verbose Output

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

### Run Tests for Specific Package

```bash
# Domain entities
go test ./domain/entities -v

# Application services
go test ./application/services -v

# Infrastructure
go test ./infrastructure/persistence -v

# HTTP handlers
go test ./interfaces/http/handlers -v

# Middleware
go test ./interfaces/http/middleware -v
```

### Generate Detailed Coverage Report

```bash
# Generate coverage profile
go test ./... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

## Test Categories

### Unit Tests

All tests are unit tests that:

- Test individual components in isolation
- Use mocks for dependencies
- Focus on business logic and edge cases
- Are fast and reliable

### Test Patterns Used

1. **Table-Driven Tests**: Most tests use table-driven patterns for comprehensive coverage
2. **Mock Objects**: Custom mock implementations for testing dependencies
3. **Test Fixtures**: Consistent test data setup and teardown
4. **Error Testing**: Comprehensive error condition testing
5. **Concurrency Testing**: Thread safety tests for shared resources

## Test Examples

### Entity Tests

Tests for the User entity cover:

- Constructor validation
- Business rule validation (`IsValid()`)
- State mutation methods
- JSON serialization compatibility

### Service Tests

Tests for UserService cover:

- CRUD operations
- Input validation
- Error handling
- Repository interaction
- Business logic edge cases

### Handler Tests

Tests for HTTP handlers cover:

- Request/response cycles
- Status code validation
- JSON marshaling/unmarshaling
- Error response formatting
- Route parameter handling

### Repository Tests

Tests for the memory repository cover:

- Data persistence operations
- Concurrency safety
- Error conditions
- Thread safety with multiple goroutines

### Middleware Tests

Tests for logging middleware cover:

- Request logging functionality
- Handler chain execution
- Response preservation
- Error propagation

## Test Data

Tests use predictable test data:

- User ID: 1, 2, 3, etc.
- Names: "John Doe", "Jane Smith", "Bob Johnson"
- Emails: "john@example.com", "jane@example.com", etc.

## Mock Objects

### MockUserRepository

Implements `UserRepository` interface with:

- Configurable failure modes
- In-memory data storage
- Error simulation capabilities

### MockUserService

Implements `UserServiceInterface` with:

- Simplified business logic
- Configurable error conditions
- Predictable behavior for testing

## Best Practices

### Test Organization

1. **One test file per source file**
2. **Group related tests in functions**
3. **Use descriptive test names**
4. **Follow AAA pattern**: Arrange, Act, Assert

### Test Naming

Tests follow the pattern: `Test[StructName]_[MethodName]`

Examples:

- `TestUser_IsValid`
- `TestUserService_CreateUser`
- `TestUserHandler_GetUsers`

### Test Structure

```go
func TestFunction_Scenario(t *testing.T) {
    tests := []struct {
        name           string
        input          InputType
        expected       ExpectedType
        expectedError  string
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Continuous Integration

Tests are designed to:

- Run quickly (all tests complete in < 5 seconds)
- Be deterministic and reliable
- Not require external dependencies
- Work in any environment

## Troubleshooting

### Common Issues

1. **Import Errors**: Ensure all dependencies are properly imported
2. **Race Conditions**: Run tests with `-race` flag to detect data races
3. **Timeout Issues**: Use context with timeouts for long-running operations

### Debug Commands

```bash
# Run with race detection
go test ./... -race

# Run specific test with verbose output
go test -v -run TestUserService_CreateUser ./application/services

# Run with CPU profiling
go test -cpuprofile=cpu.prof ./...
```

## Contributing

When adding new features:

1. **Write tests first** (TDD approach)
2. **Maintain 100% coverage** for new code
3. **Follow existing patterns** for consistency
4. **Update this README** if test structure changes

### Test Checklist

- [ ] Unit tests for all public methods
- [ ] Error condition testing
- [ ] Edge case coverage
- [ ] Mock object usage for dependencies
- [ ] Table-driven tests where appropriate
- [ ] Descriptive test names and documentation
