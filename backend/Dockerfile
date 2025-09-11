# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25.0
FROM golang:${GO_VERSION}-alpine AS builder

# Install git for go get if needed, and ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY --link go.mod go.mod
COPY --link go.sum go.sum

# Download dependencies (with cache mounts for speed)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy the rest of the source code
COPY --link . .

# Build the Go binary (static build)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app

# --- Final minimal image ---
FROM alpine:latest AS final

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy the built binary from builder
COPY --from=builder /app/app /app/app
# COPY --from=builder /app/bac.sql /app/bac.sql

# Use non-root user
USER appuser

# Expose the port (as per main.go:8080)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/app/app"]
