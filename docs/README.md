# Secretly Documentation

Welcome to the Secretly documentation, a modern web application for securely and centrally managing environment variables.

## 📚 Deployment Guides

### [Docker Compose](./docker-compose.md)
Complete guide for deploying Secretly using Docker Compose:
- Basic and advanced configuration
- Configurations for development and production
- Integration with external databases
- Monitoring and health checks
- Common troubleshooting
- CI/CD integration

### [Kubernetes](./kubernetes.md)
Complete guide for deploying Secretly on Kubernetes:
- Basic YAML configuration
- Configurations for different environments
- Horizontal Pod Autoscaler (HPA)
- Monitoring with Prometheus
- Network Policies and security
- Integration with ArgoCD and Flux

## 🚀 Quick Start

### With Docker
```bash
docker run -d \
  --name secretly \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  rodrwan/secretly:latest
```

### With Docker Compose
```bash
# Clone the repository
git clone https://github.com/rodrwan/secretly.git
cd secretly

# Start the service
docker-compose up -d
```

### With Kubernetes
```bash
# Apply basic configuration
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/pvc.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `PORT` | Port where the server runs | `8080` |

### API Endpoints

- `GET /api/v1/env` - Get all environment variables
- `POST /api/v1/env` - Create/update environment variables
- `GET /api/v1/env/{id}` - Get a specific environment
- `PUT /api/v1/env/{id}` - Update a specific environment
- `DELETE /api/v1/env/{id}` - Delete a specific environment

## 🛠️ Development

### Prerequisites
- Go 1.24 or higher
- Docker and Docker Compose
- Make

### Development Commands
```bash
make help        # Show available commands
make dev         # Run locally
make docker-dev  # Run with Docker
make generate    # Generate template files
```

## 🔒 Security

### Production Recommendations

1. **Use named volumes** instead of bind mounts
2. **Configure secrets** for sensitive credentials
3. **Use custom networks** to isolate services
4. **Implement authentication** if necessary
5. **Configure automatic backups** of data
6. **Use HTTPS** with valid certificates
7. **Implement rate limiting** in ingress

### Secure Configuration Example

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  secretly:
    image: rodrwan/secretly:latest
    ports:
      - "127.0.0.1:8080:8080"  # Local access only
    volumes:
      - secretly-data:/app/data  # Named volume
    environment:
      - PORT=8080
    secrets:
      - db_password
    networks:
      - secretly-network

volumes:
  secretly-data:
    driver: local

secrets:
  db_password:
    file: ./secrets/db_password.txt

networks:
  secretly-network:
    driver: bridge
```

## 📊 Monitoring

### Health Checks
The service includes automatic health checks that verify:
- Port 8080 availability
- HTTP response from main endpoint
- Container status

### Metrics
- Health check endpoint: `GET /api/v1/env`
- Structured logs in JSON format
- Prometheus metrics (optional)

## 🔄 CI/CD

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
      - name: Deploy
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

## 🆘 Troubleshooting

### Common Issues

1. **Container doesn't start**
   - Check logs: `docker-compose logs secretly`
   - Check data directory permissions
   - Recreate container: `docker-compose up -d --force-recreate`

2. **Data doesn't persist**
   - Check volume configuration
   - Check write permissions
   - Recreate volume: `docker-compose down -v && docker-compose up -d`

3. **API doesn't respond**
   - Check health checks
   - Check network configuration
   - Check service logs

### Diagnostic Commands

```bash
# Docker Compose
docker-compose ps
docker-compose logs -f secretly
docker-compose exec secretly sh

# Kubernetes
kubectl get pods -n secretly
kubectl logs -f deployment/secretly -n secretly
kubectl describe pod <pod-name> -n secretly
```

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/rodrwan/secretly/issues)
- **Documentation**: [Main README](../README.md)
- **Examples**: [Usage examples](../examples/)

## 📄 License

MIT License - see [LICENSE](../LICENSE) file for details. 