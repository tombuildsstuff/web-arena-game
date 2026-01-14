# Stage 1: Build the client
FROM oven/bun:1 AS client-builder

WORKDIR /app/client

# Copy client package files
COPY client/package.json client/bun.lock* ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy client source
COPY client/ ./

# Build the client
RUN bun run build

# Stage 2: Build the server
FROM golang:1.25-alpine AS server-builder

WORKDIR /app/server

# Install build dependencies
RUN apk add --no-cache make

# Copy go mod files first for caching
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod download

# Copy server source
COPY server/ ./

# Copy built client to static directory
COPY --from=client-builder /app/client/dist ./static/

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /arena-server main.go

# Stage 3: Runtime
FROM alpine:3.20

# Add CA certificates for HTTPS and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=server-builder /arena-server .

# Use non-root user
USER appuser

# Expose the server port
EXPOSE 3000

# Run the server
ENTRYPOINT ["./arena-server"]
