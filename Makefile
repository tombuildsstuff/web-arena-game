.PHONY: tools build run run-dev clean test help

# Default target
all: build

# Install all development tools and dependencies
tools:
	@echo "Installing tools for server and client..."
	@$(MAKE) -C server tools
	@$(MAKE) -C client tools
	@echo "✓ All tools installed successfully"

# Build server (client is automatically embedded as a dependency)
build:
	@$(MAKE) -C server build

# Run the server (client embedded inside)
run:
	@$(MAKE) -C server run

# Run in development mode (hot reload for both server and client)
run-dev:
	@echo "Starting development servers..."
	@echo "  Server API: http://localhost:3000"
	@echo "  Client Dev: http://localhost:5173 (with hot reload)"
	@echo ""
	@echo "Note: Use Ctrl+C to stop both servers"
	@make -j2 run-dev-server run-dev-client

# Run server in dev mode (go run, shows fallback page since no embedded content)
run-dev-server:
	@cd server && go run main.go

# Run client dev server with hot reload
run-dev-client:
	@$(MAKE) -C client run

# Run tests
test:
	@$(MAKE) -C server test

# Clean all build artifacts
clean:
	@$(MAKE) -C server clean
	@$(MAKE) -C client clean

# Show help
help:
	@echo "Arena Game - Makefile Commands"
	@echo ""
	@echo "Quick Start:"
	@echo "  make tools    - Install all dependencies"
	@echo "  make build    - Build server (client auto-embedded)"
	@echo "  make run      - Run the server"
	@echo ""
	@echo "Main Commands:"
	@echo "  make tools    - Install Go tools and Bun dependencies"
	@echo "  make build    - Build server with embedded client"
	@echo "  make run      - Build and run server (http://localhost:3000)"
	@echo "  make run-dev  - Run dev servers with hot reload"
	@echo "  make clean    - Remove all build artifacts"
	@echo "  make test     - Run server tests"
	@echo ""
	@echo "Development Mode (Hot Reload):"
	@echo "  make run-dev  - Run both servers with hot reload enabled"
	@echo "    Client: http://localhost:5173 (Vite dev server)"
	@echo "    Server: http://localhost:3000 (Go server)"
	@echo ""
	@echo "Production Mode (Single Binary):"
	@echo "  make build    - Build server/bin/arena-server"
	@echo "  make run      - Run the built binary"
	@echo "  Direct: ./server/bin/arena-server"
	@echo ""
	@echo "Note: The client is ALWAYS embedded in the server binary."
	@echo "      Use 'make run' for production-like testing."
	@echo "      Use 'make run-dev' for development with hot reload."
