#!/bin/bash

# Kubernetes deployment script for Secretly
# Usage: ./deploy.sh [dev|prod]

set -e

NAMESPACE="secretly"
ENVIRONMENT=${1:-dev}

echo "🚀 Deploying Secretly on Kubernetes (environment: $ENVIRONMENT)"

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed"
    exit 1
fi

# Create namespace
echo "📦 Creating namespace..."
kubectl apply -f namespace.yaml

# Apply ConfigMap
echo "⚙️  Applying ConfigMap..."
kubectl apply -f configmap.yaml

# Apply PVC
echo "💾 Applying PersistentVolumeClaim..."
kubectl apply -f pvc.yaml

# Apply Deployment
echo "🔄 Applying Deployment..."
kubectl apply -f deployment.yaml

# Apply Service
echo "🌐 Applying Service..."
kubectl apply -f service.yaml

# Apply Ingress if production
if [ "$ENVIRONMENT" = "prod" ]; then
    echo "🔗 Applying Ingress..."
    kubectl apply -f ingress.yaml
fi

# Wait for pods to be ready
echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=secretly -n $NAMESPACE --timeout=300s

# Check deployment
echo "✅ Checking deployment..."
kubectl get all -n $NAMESPACE

echo "🎉 Deployment completed!"
echo ""
echo "📊 Deployment status:"
kubectl get pods -n $NAMESPACE
echo ""
echo "🌐 To access the service:"
if [ "$ENVIRONMENT" = "prod" ]; then
    echo "   URL: https://secretly.yourdomain.com"
else
    echo "   Port-forward: kubectl port-forward service/secretly-service 8080:80 -n $NAMESPACE"
    echo "   Local URL: http://localhost:8080"
fi
echo ""
echo "📝 Logs: kubectl logs -f deployment/secretly -n $NAMESPACE"