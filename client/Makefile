.PHONY: tools build run clean preview

# Default target
all: build

# Install Bun dependencies
tools:
	@echo "Installing dependencies with Bun..."
	bun install
	@echo "Dependencies installed successfully"

# Build the client for production
build:
	@echo "Building client for production..."
	bun run build
	@echo "Build complete: dist/"

# Run the development server
run:
	@echo "Starting development server on :5173..."
	bun run dev

# Preview production build locally
preview:
	@echo "Starting preview server..."
	bun run preview

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist/
	rm -rf node_modules/
	@echo "Clean complete"

# Update dependencies
update:
	@echo "Updating dependencies..."
	bun update
	@echo "Dependencies updated"

# Reinstall dependencies
reinstall: clean tools
	@echo "Dependencies reinstalled"
