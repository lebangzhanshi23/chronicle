# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o chronicle main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/chronicle .
COPY --from=builder /app/frontend/dist ./frontend/dist
COPY --from=builder /app/templates ./templates

# Create data directory
RUN mkdir -p data

# Expose port
EXPOSE 8080

# Run the application (default to server)
CMD ["./chronicle", "server"]
