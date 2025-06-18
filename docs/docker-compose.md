# Despliegue con Docker Compose

Esta guía te ayudará a desplegar Secretly usando Docker Compose en diferentes ambientes de desarrollo.

## Configuración Básica

### 1. Archivo docker-compose.yml

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
      - ENV_PATH=/app/data/.env
      - BASE_PATH=/api/v1
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

### 2. Despliegue

```bash
# Crear el directorio de datos
mkdir -p data

# Iniciar el servicio
docker-compose up -d

# Verificar el estado
docker-compose ps

# Ver logs
docker-compose logs -f secretly
```

## Configuraciones Avanzadas

### Desarrollo Local

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
      - ./internal:/app/internal  # Para desarrollo con hot-reload
    environment:
      - PORT=8080
      - ENV_PATH=/app/data/.env
      - BASE_PATH=/api/v1
      - DEBUG=true
    restart: unless-stopped
    networks:
      - secretly-network

networks:
  secretly-network:
    driver: bridge
```

### Producción

```yaml
version: '3.8'

services:
  secretly:
    image: rodrwan/secretly:latest
    container_name: secretly-prod
    ports:
      - "9000:8080"  # Puerto externo diferente
    volumes:
      - secretly-data:/app/data
    environment:
      - PORT=8080
      - ENV_PATH=/app/data/.env
      - BASE_PATH=/api/v1
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

### Con Base de Datos Externa

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
      - ENV_PATH=/app/data/.env
      - BASE_PATH=/api/v1
      - DATABASE_URL=postgresql://user:password@postgres:5432/secretly
    depends_on:
      - postgres
    restart: unless-stopped
    networks:
      - secretly-network

  postgres:
    image: postgres:15-alpine
    container_name: secretly-postgres
    environment:
      - POSTGRES_DB=secretly
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - secretly-network

volumes:
  postgres-data:
    driver: local

networks:
  secretly-network:
    driver: bridge
```

## Variables de Entorno

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto donde se ejecuta el servidor | `8080` |
| `ENV_PATH` | Ruta al archivo .env | `/app/data/.env` |
| `BASE_PATH` | Ruta base de la API | `/api/v1` |
| `DEBUG` | Modo debug (desarrollo) | `false` |
| `DATABASE_URL` | URL de conexión a base de datos | - |

## Comandos Útiles

```bash
# Iniciar servicios
docker-compose up -d

# Detener servicios
docker-compose down

# Reiniciar un servicio específico
docker-compose restart secretly

# Ver logs en tiempo real
docker-compose logs -f secretly

# Ejecutar comandos dentro del contenedor
docker-compose exec secretly sh

# Hacer backup de los datos
docker-compose exec secretly tar -czf /tmp/backup.tar.gz /app/data
docker cp secretly:/tmp/backup.tar.gz ./backup.tar.gz

# Restaurar backup
docker cp ./backup.tar.gz secretly:/tmp/backup.tar.gz
docker-compose exec secretly tar -xzf /tmp/backup.tar.gz -C /
```

## Monitoreo y Health Checks

El servicio incluye health checks automáticos que verifican:

- Disponibilidad del puerto 8080
- Respuesta HTTP del endpoint principal
- Estado del contenedor

### Verificar Estado

```bash
# Verificar health check
docker-compose ps

# Ver logs de health check
docker inspect secretly | grep -A 10 "Health"

# Verificar endpoint manualmente
curl http://localhost:8080/api/v1/env
```

## Troubleshooting

### Problema: Contenedor no inicia

```bash
# Verificar logs
docker-compose logs secretly

# Verificar permisos del directorio data
ls -la data/

# Recrear contenedor
docker-compose down
docker-compose up -d --force-recreate
```

### Problema: Datos no persisten

```bash
# Verificar volumen
docker volume ls

# Verificar montaje
docker inspect secretly | grep -A 10 "Mounts"

# Recrear volumen
docker-compose down -v
docker-compose up -d
```

### Problema: Puerto ocupado

```bash
# Cambiar puerto en docker-compose.yml
ports:
  - "8081:8080"  # Puerto externo 8081

# O verificar qué usa el puerto
lsof -i :8080
```

## Seguridad

### Recomendaciones para Producción

1. **Usar volúmenes nombrados** en lugar de bind mounts
2. **Configurar secrets** para credenciales sensibles
3. **Usar redes personalizadas** para aislar servicios
4. **Implementar autenticación** si es necesario
5. **Configurar backups automáticos** de los datos

### Ejemplo con Secrets

```yaml
version: '3.8'

services:
  secretly:
    image: rodrwan/secretly:latest
    secrets:
      - db_password
    environment:
      - DATABASE_PASSWORD_FILE=/run/secrets/db_password
    # ... resto de configuración

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

## Integración con CI/CD

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