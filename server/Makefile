.PHONY: tools build run clean test fmt lint deps

# Default target
all: build

# Install Go development tools
tools:
	@echo "Installing Go tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Build the server (always embeds the client)
build:
	@echo "Building client..."
	@$(MAKE) -C ../client build
	@echo "Copying client to static directory..."
	@rm -rf static/
	@mkdir -p static
	@cp -r ../client/dist/* static/
	@echo "Building server with embedded frontend..."
	@go build -o bin/arena-server main.go
	@echo "✓ Build complete: bin/arena-server"
	@echo "  Run with: make run or ./bin/arena-server"

# Run the server binary
run: build
	@echo "Starting server on :3000..."
	@echo "  Web UI: http://localhost:3000"
	@echo "  WebSocket: ws://localhost:3000/ws"
	@./bin/arena-server

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf static/
	@echo "Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated"
