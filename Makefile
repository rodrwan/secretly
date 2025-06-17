.PHONY: build run test clean docker-build docker-up docker-down generate-templ help

# Variables
BINARY_NAME=secretly
DOCKER_COMPOSE=docker-compose

# Default target
.DEFAULT_GOAL := help

# Help command
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run          - Run the application locally"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make generate-templ - Generate templ files"

# Build the application
build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) ./cmd/server

# Run the application
run: build
	@echo "Running application..."
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean

# Docker commands
docker-build:
	@echo "Building Docker image..."
	$(DOCKER_COMPOSE) build

docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

# Generate templ files
generate-templ:
	@echo "Generating templ files..."
	templ generate ./internal/web/templates/

# Development workflow
dev: generate-templ build run

# Docker development workflow
docker-dev: docker-build docker-up

# Publish the image
publish:
	@echo "Publishing image..."
	./scripts/publish.sh