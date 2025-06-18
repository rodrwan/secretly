# Documentación de Secretly

Bienvenido a la documentación de Secretly, una aplicación web moderna para gestionar variables de entorno de forma segura y centralizada.

## 📚 Guías de Despliegue

### [Docker Compose](./docker-compose.md)
Guía completa para desplegar Secretly usando Docker Compose:
- Configuración básica y avanzada
- Configuraciones para desarrollo y producción
- Integración con bases de datos externas
- Monitoreo y health checks
- Troubleshooting común
- Integración con CI/CD

### [Kubernetes](./kubernetes.md)
Guía completa para desplegar Secretly en Kubernetes:
- Configuración básica con YAML
- Configuraciones para diferentes ambientes
- Horizontal Pod Autoscaler (HPA)
- Monitoreo con Prometheus
- Network Policies y seguridad
- Integración con ArgoCD y Flux

## 🚀 Inicio Rápido

### Con Docker
```bash
docker run -d \
  --name secretly \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  rodrwan/secretly:latest
```

### Con Docker Compose
```bash
# Clonar el repositorio
git clone https://github.com/rodrwan/secretly.git
cd secretly

# Iniciar el servicio
docker-compose up -d
```

### Con Kubernetes
```bash
# Aplicar configuración básica
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/pvc.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## 🔧 Configuración

### Variables de Entorno

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto donde se ejecuta el servidor | `8080` |

### API Endpoints

- `GET /api/v1/env` - Obtener todas las variables de entorno
- `POST /api/v1/env` - Crear/actualizar variables de entorno
- `GET /api/v1/env/{id}` - Obtener un entorno específico
- `PUT /api/v1/env/{id}` - Actualizar un entorno específico
- `DELETE /api/v1/env/{id}` - Eliminar un entorno específico

## 🛠️ Desarrollo

### Prerrequisitos
- Go 1.24 o superior
- Docker y Docker Compose
- Make

### Comandos de Desarrollo
```bash
make help        # Mostrar comandos disponibles
make dev         # Ejecutar localmente
make docker-dev  # Ejecutar con Docker
make generate    # Generar archivos de templates
```

## 🔒 Seguridad

### Recomendaciones para Producción

1. **Usar volúmenes nombrados** en lugar de bind mounts
2. **Configurar secrets** para credenciales sensibles
3. **Usar redes personalizadas** para aislar servicios
4. **Implementar autenticación** si es necesario
5. **Configurar backups automáticos** de los datos
6. **Usar HTTPS** con certificados válidos
7. **Implementar rate limiting** en el ingress

### Ejemplo de Configuración Segura

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  secretly:
    image: rodrwan/secretly:latest
    ports:
      - "127.0.0.1:8080:8080"  # Solo acceso local
    volumes:
      - secretly-data:/app/data  # Volumen nombrado
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

## 📊 Monitoreo

### Health Checks
El servicio incluye health checks automáticos que verifican:
- Disponibilidad del puerto 8080
- Respuesta HTTP del endpoint principal
- Estado del contenedor

### Métricas
- Endpoint de health check: `GET /api/v1/env`
- Logs estructurados en formato JSON
- Métricas de Prometheus (opcional)

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

### Problemas Comunes

1. **Contenedor no inicia**
   - Verificar logs: `docker-compose logs secretly`
   - Verificar permisos del directorio data
   - Recrear contenedor: `docker-compose up -d --force-recreate`

2. **Datos no persisten**
   - Verificar configuración de volúmenes
   - Verificar permisos de escritura
   - Recrear volumen: `docker-compose down -v && docker-compose up -d`

3. **API no responde**
   - Verificar health checks
   - Verificar configuración de red
   - Verificar logs del servicio

### Comandos de Diagnóstico

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

## 📞 Soporte

- **Issues**: [GitHub Issues](https://github.com/rodrwan/secretly/issues)
- **Documentación**: [README principal](../README.md)
- **Ejemplos**: [Ejemplos de uso](../examples/)

## 📄 Licencia

MIT License - ver archivo [LICENSE](../LICENSE) para más detalles. 