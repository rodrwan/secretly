# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate ./internal/web/templates/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/secretly ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/secretly .

# Create directory for .env file
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./secretly"]