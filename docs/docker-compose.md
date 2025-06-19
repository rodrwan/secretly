# Docker Compose Deployment

This guide will help you deploy Secretly using Docker Compose in different development environments.

## Basic Configuration

### 1. docker-compose.yml File

```yaml
version: '3.8'

services:
  secretly:
    image: rodrwan/secretly:latest
    container_name: secretly
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - PORT=8080
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - secretly-network

networks:
  secretly-network:
    driver: bridge
```

### 2. Deployment

```bash
# Create data directory
mkdir -p data

# Start the service
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f secretly
```

## Advanced Configurations

### Local Development

```yaml
version: '3.8'

services:
  secretly:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: secretly-dev
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./internal:/app/internal  # For development with hot-reload
    environment:
      - PORT=8080
      - DEBUG=true
    restart: unless-stopped
    networks:
      - secretly-network

networks:
  secretly-network:
    driver: bridge
```

### Production

```yaml
version: '3.8'

services:
  secretly:
    image: rodrwan/secretly:latest
    container_name: secretly-prod
    ports:
      - "9000:8080"  # Different external port
    volumes:
      - secretly-data:/app/data
    environment:
      - PORT=8080
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - secretly-network
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.secretly.rule=Host(`secretly.yourdomain.com`)"
      - "traefik.http.routers.secretly.tls=true"

volumes:
  secretly-data:
    driver: local

networks:
  secretly-network:
    driver: bridge
```

## Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `PORT` | Port where the server runs | `8080` |

## Useful Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Restart a specific service
docker-compose restart secretly

# View logs in real-time
docker-compose logs -f secretly

# Execute commands inside the container
docker-compose exec secretly sh

# Backup data
docker-compose exec secretly tar -czf /tmp/backup.tar.gz /app/data
docker cp secretly:/tmp/backup.tar.gz ./backup.tar.gz

# Restore backup
docker cp ./backup.tar.gz secretly:/tmp/backup.tar.gz
docker-compose exec secretly tar -xzf /tmp/backup.tar.gz -C /
```

## Monitoring and Health Checks

The service includes automatic health checks that verify:

- Port 8080 availability
- HTTP response from main endpoint
- Container status

### Check Status

```bash
# Check health status
docker-compose ps

# View health check logs
docker inspect secretly | grep -A 10 "Health"

# Manually check endpoint
curl http://localhost:8080/api/v1/env
```

## Troubleshooting

### Issue: Container doesn't start

```bash
# Check logs
docker-compose logs secretly

# Check data directory permissions
ls -la data/

# Recreate container
docker-compose down
docker-compose up -d --force-recreate
```

### Issue: Data doesn't persist

```bash
# Check volume
docker volume ls

# Check mount
docker inspect secretly | grep -A 10 "Mounts"

# Recreate volume
docker-compose down -v
docker-compose up -d
```

### Issue: Port is occupied

```bash
# Change port in docker-compose.yml
ports:
  - "8081:8080"  # External port 8081

# Or check what's using the port
lsof -i :8080
```

## Security

### Production Recommendations

1. **Use named volumes** instead of bind mounts
2. **Configure secrets** for sensitive credentials
3. **Use custom networks** to isolate services
4. **Implement authentication** if necessary
5. **Configure automatic backups** of data

### Example with Secrets

```yaml
version: '3.8'

services:
  secretly:
    image: rodrwan/secretly:latest
    secrets:
      - db_password
    environment:
      - DATABASE_PASSWORD_FILE=/run/secrets/db_password
    # ... rest of configuration

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Deploy Secretly

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to server
        run: |
          ssh user@server "cd /path/to/secretly && git pull && docker-compose up -d"
```

### GitLab CI

```yaml
deploy:
  stage: deploy
  script:
    - ssh user@server "cd /path/to/secretly && git pull && docker-compose up -d"
  only:
    - main
``` 