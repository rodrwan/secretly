# Secretly 🔐

Secretly is a web application that allows you to create and manage a .env file through a simple web interface.

This service is built in Go and exposes a dashboard where you can securely manage your .env content.

It also exposes an HTTP client that allows you to consume this env remotely, simply and securely.

## Development

### Prerequisites
- Go 1.24 or later
- Docker and Docker Compose (for containerized deployment)
- Make (for using Makefile commands)

### Quick Start

#### Using Make
The project includes a Makefile with common commands:

```bash
# Show available commands
make help

# Development workflow (build and run locally)
make dev

# Docker development workflow
make docker-dev

# Generate templ files
make generate-templ

# Run tests
make test
```

#### Manual Setup
1. Clone the repository:
```bash
git clone https://github.com/rodrwan/secretly.git
cd secretly
```

2. Start the application:
```bash
docker-compose up -d
```

The application will be available at http://localhost:8080

### Environment Variables
The following environment variables can be configured in the `docker-compose.yml` file:

- `PORT`: The port where the application will run (default: 8080)
- `ENV_PATH`: Path to the .env file (default: /app/data/.env)
- `BASE_PATH`: Base path for the API (default: /api/v1)

### Data Persistence
The .env file is stored in the `./data` directory and is persisted through a Docker volume.

### Stopping the Application
```bash
# Using Make
make docker-down

# Or using Docker Compose directly
docker-compose down
```


