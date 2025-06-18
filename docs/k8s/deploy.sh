#!/bin/bash

# Script de despliegue para Secretly en Kubernetes
# Uso: ./deploy.sh [dev|prod]

set -e

NAMESPACE="secretly"
ENVIRONMENT=${1:-dev}

echo "🚀 Desplegando Secretly en Kubernetes (ambiente: $ENVIRONMENT)"

# Verificar que kubectl esté disponible
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl no está instalado"
    exit 1
fi

# Crear namespace
echo "📦 Creando namespace..."
kubectl apply -f namespace.yaml

# Aplicar ConfigMap
echo "⚙️  Aplicando ConfigMap..."
kubectl apply -f configmap.yaml

# Aplicar PVC
echo "💾 Aplicando PersistentVolumeClaim..."
kubectl apply -f pvc.yaml

# Aplicar Deployment
echo "🔄 Aplicando Deployment..."
kubectl apply -f deployment.yaml

# Aplicar Service
echo "🌐 Aplicando Service..."
kubectl apply -f service.yaml

# Aplicar Ingress si es producción
if [ "$ENVIRONMENT" = "prod" ]; then
    echo "🔗 Aplicando Ingress..."
    kubectl apply -f ingress.yaml
fi

# Esperar a que los pods estén listos
echo "⏳ Esperando a que los pods estén listos..."
kubectl wait --for=condition=ready pod -l app=secretly -n $NAMESPACE --timeout=300s

# Verificar el despliegue
echo "✅ Verificando despliegue..."
kubectl get all -n $NAMESPACE

echo "🎉 ¡Despliegue completado!"
echo ""
echo "📊 Estado del despliegue:"
kubectl get pods -n $NAMESPACE
echo ""
echo "🌐 Para acceder al servicio:"
if [ "$ENVIRONMENT" = "prod" ]; then
    echo "   URL: https://secretly.yourdomain.com"
else
    echo "   Port-forward: kubectl port-forward service/secretly-service 8080:80 -n $NAMESPACE"
    echo "   URL local: http://localhost:8080"
fi
echo ""
echo "📝 Logs: kubectl logs -f deployment/secretly -n $NAMESPACE"