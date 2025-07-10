# justfile for hanacaraka Go project

# Load environment variables
set dotenv-load := true

# Get PORT from environment, default to 8080 if not set
port := env_var_or_default('PORT', '8080')

# Default recipe to display available commands
default:
    @just --list

# Build the Go binary
build:
    @echo "Building hanacaraka..."
    go build -o bin/hanacaraka .

# Build and run the binary
start: build
    @echo "Starting hanacaraka server..."
    ./bin/hanacaraka

# Run the server with hot reload using air
dev:
    @echo "Starting development server with hot reload..."
    @echo "Using port: {{port}}"
    @echo "Checking if port {{port}} is in use..."
    @if lsof -ti:{{port}} > /dev/null 2>&1; then \
        echo "Port {{port}} is in use. Killing existing process..."; \
        lsof -ti:{{port}} | xargs kill -9; \
        sleep 1; \
    fi
    @if ! command -v air > /dev/null 2>&1; then \
        echo "Installing air for hot reload..."; \
        go install github.com/air-verse/air@latest; \
    fi
    @if command -v air > /dev/null 2>&1; then \
        air; \
    else \
        echo "Using air from GOPATH/bin..."; \
        $(go env GOPATH)/bin/air; \
    fi

# Run all tests in the project
test:
    @echo "Running all tests..."
    go test -v ./...

# Run tests with coverage
test-coverage:
    @echo "Running tests with coverage..."
    go test -v -cover ./...

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf bin/
    go clean

# Format Go code
fmt:
    @echo "Formatting Go code..."
    go fmt ./...

# Run Go linter
lint:
    @echo "Running Go linter..."
    @if ! command -v golangci-lint > /dev/null 2>&1; then \
        echo "Installing golangci-lint..."; \
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
    fi
    golangci-lint run

# Download and tidy Go modules
deps:
    @echo "Downloading and tidying Go modules..."
    go mod download
    go mod tidy

# Run security check
security:
    @echo "Running security check..."
    @if ! command -v gosec > /dev/null 2>&1; then \
        echo "Installing gosec..."; \
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
    fi
    gosec ./...

# Kill any process using the configured port
kill-port:
    @echo "Killing any process using port {{port}}..."
    @if lsof -ti:{{port}} > /dev/null 2>&1; then \
        echo "Found process using port {{port}}. Killing..."; \
        lsof -ti:{{port}} | xargs kill -9; \
        echo "Process killed."; \
    else \
        echo "No process found using port {{port}}."; \
    fi

# Show current environment configuration
env-info:
    @echo "Current environment configuration:"
    @echo "PORT: {{port}}"
    @echo "ENV: $(printenv ENV || echo 'not set')"
    @echo "HOST: $(printenv HOST || echo 'not set')"

# Load development environment
dev-env:
    @echo "Loading development environment..."
    @cp .env.development .env
    @echo "Development environment loaded. Run 'just dev' to start."

# Load production environment
prod-env:
    @echo "Loading production environment..."
    @cp .env.production .env
    @echo "Production environment loaded."
