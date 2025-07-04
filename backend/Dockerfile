# Multi-stage Dockerfile for Backend
FROM golang:1.22-alpine AS base

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod first
COPY go.mod ./

# Copy go.sum if it exists (optional)
COPY go.su[m] ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Development stage
FROM base AS development
# Install Air for hot reloading (use compatible version)
RUN go install github.com/cosmtrek/air@v1.49.0
EXPOSE 8080
CMD ["air"]

# Testing stage
FROM base AS testing
# Run tests
RUN go test -v ./...
EXPOSE 8080
CMD ["go", "test", "-v", "./..."]

# Builder stage
FROM base AS builder
# Build the application
RUN go build -o main .

# Production stage
FROM alpine:latest AS production
WORKDIR /app
# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata
# Copy the binary from builder stage
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]