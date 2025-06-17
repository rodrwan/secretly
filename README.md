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

## Client Integration

You can integrate the Secretly client into your Go applications to remotely access environment variables. Here's how to do it:

### Installation

Add Secretly to your Go project:

```bash
go get github.com/rodrwan/secretly
```

### Usage

Here's a simple example of how to use the Secretly client in your Go application:

```go
package main

import (
    "fmt"
    "log"

    "github.com/rodrwan/secretly/pkg/secretly"
)

func main() {
    // Create a new client
    client := secretly.NewClient("http://localhost:8080")

    // Get all environment variables
    vars, err := client.GetAll()
    if err != nil {
        log.Fatal(err)
    }

    // Access specific variables
    fmt.Println("Database URL:", vars["DATABASE_URL"])
    fmt.Println("API Key:", vars["API_KEY"])

    // Get a specific variable
    dbURL, err := client.Get("DATABASE_URL")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Database URL:", dbURL)
}
```

### Client Configuration

The client can be configured with different options:

```go
client := secretly.NewClient(
    "http://localhost:8080",
    secretly.WithBasePath("/api/v1"),
    secretly.WithTimeout(5 * time.Second),
)
```

### Error Handling

The client provides detailed error information:

```go
vars, err := client.GetAll()
if err != nil {
    switch {
    case secretly.IsNotFound(err):
        fmt.Println("Environment variables not found")
    case secretly.IsUnauthorized(err):
        fmt.Println("Unauthorized access")
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

### Best Practices

1. **Caching**: Consider implementing caching for frequently accessed variables:
```go
type CachedClient struct {
    client    *secretly.Client
    cache     map[string]string
    cacheTime time.Time
}

func (c *CachedClient) Get(key string) (string, error) {
    if time.Since(c.cacheTime) > 5*time.Minute {
        vars, err := c.client.GetAll()
        if err != nil {
            return "", err
        }
        c.cache = vars
        c.cacheTime = time.Now()
    }
    return c.cache[key], nil
}
```

2. **Environment Fallback**: Always provide fallback values for critical variables:
```go
func getConfig(client *secretly.Client) (*Config, error) {
    vars, err := client.GetAll()
    if err != nil {
        return &Config{
            DatabaseURL: os.Getenv("DATABASE_URL"),
            APIKey:     os.Getenv("API_KEY"),
        }, nil
    }
    return &Config{
        DatabaseURL: vars["DATABASE_URL"],
        APIKey:     vars["API_KEY"],
    }, nil
}
```

3. **Health Checks**: Implement health checks for the Secretly service:
```go
func isSecretlyHealthy(client *secretly.Client) bool {
    _, err := client.GetAll()
    return err == nil
}
```


