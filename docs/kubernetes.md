# Despliegue en Kubernetes

Esta guía te ayudará a desplegar Secretly en un cluster de Kubernetes, incluyendo configuraciones para diferentes ambientes.

## Configuración Básica

### 1. Namespace

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: secretly
  labels:
    name: secretly
```

### 2. ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: secretly-config
  namespace: secretly
data:
  PORT: "8080"
  ENV_PATH: "/app/data/.env"
  BASE_PATH: "/api/v1"
```

### 3. Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secretly
  namespace: secretly
  labels:
    app: secretly
spec:
  replicas: 1
  selector:
    matchLabels:
      app: secretly
  template:
    metadata:
      labels:
        app: secretly
    spec:
      containers:
      - name: secretly
        image: rodrwan/secretly:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: PORT
          valueFrom:
            configMapKeyRef:
              name: secretly-config
              key: PORT
        - name: ENV_PATH
          valueFrom:
            configMapKeyRef:
              name: secretly-config
              key: ENV_PATH
        - name: BASE_PATH
          valueFrom:
            configMapKeyRef:
              name: secretly-config
              key: BASE_PATH
        volumeMounts:
        - name: data-volume
          mountPath: /app/data
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        livenessProbe:
          httpGet:
            path: /api/v1/env
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/env
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: data-volume
        persistentVolumeClaim:
          claimName: secretly-pvc
```

### 4. PersistentVolumeClaim

```yaml
# pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: secretly-pvc
  namespace: secretly
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: standard  # Ajustar según tu cluster
```

### 5. Service

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: secretly-service
  namespace: secretly
  labels:
    app: secretly
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: secretly
```

### 6. Ingress (Opcional)

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: secretly-ingress
  namespace: secretly
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - secretly.yourdomain.com
    secretName: secretly-tls
  rules:
  - host: secretly.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: secretly-service
            port:
              number: 80
```

## Despliegue

```bash
# Crear namespace
kubectl apply -f namespace.yaml

# Aplicar configuración
kubectl apply -f configmap.yaml
kubectl apply -f pvc.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Aplicar ingress (opcional)
kubectl apply -f ingress.yaml

# Verificar el despliegue
kubectl get all -n secretly
kubectl get pvc -n secretly
```

## Configuraciones Avanzadas

### Desarrollo

```yaml
# deployment-dev.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secretly-dev
  namespace: secretly
spec:
  replicas: 1
  selector:
    matchLabels:
      app: secretly-dev
  template:
    metadata:
      labels:
        app: secretly-dev
    spec:
      containers:
      - name: secretly
        image: rodrwan/secretly:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: ENV_PATH
          value: "/app/data/.env"
        - name: BASE_PATH
          value: "/api/v1"
        - name: DEBUG
          value: "true"
        volumeMounts:
        - name: data-volume
          mountPath: /app/data
        resources:
          requests:
            memory: "32Mi"
            cpu: "25m"
          limits:
            memory: "64Mi"
            cpu: "50m"
      volumes:
      - name: data-volume
        emptyDir: {}
```

### Producción con HPA

```yaml
# deployment-prod.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secretly-prod
  namespace: secretly
spec:
  replicas: 3
  selector:
    matchLabels:
      app: secretly-prod
  template:
    metadata:
      labels:
        app: secretly-prod
    spec:
      containers:
      - name: secretly
        image: rodrwan/secretly:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: ENV_PATH
          value: "/app/data/.env"
        - name: BASE_PATH
          value: "/api/v1"
        volumeMounts:
        - name: data-volume
          mountPath: /app/data
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /api/v1/env
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /api/v1/env
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
      volumes:
      - name: data-volume
        persistentVolumeClaim:
          claimName: secretly-pvc-prod
```

### Horizontal Pod Autoscaler

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: secretly-hpa
  namespace: secretly
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: secretly-prod
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Con Base de Datos Externa

```yaml
# deployment-with-db.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secretly
  namespace: secretly
spec:
  replicas: 1
  selector:
    matchLabels:
      app: secretly
  template:
    metadata:
      labels:
        app: secretly
    spec:
      containers:
      - name: secretly
        image: rodrwan/secretly:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: ENV_PATH
          value: "/app/data/.env"
        - name: BASE_PATH
          value: "/api/v1"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        volumeMounts:
        - name: data-volume
          mountPath: /app/data
      volumes:
      - name: data-volume
        persistentVolumeClaim:
          claimName: secretly-pvc
```

### Secret para Base de Datos

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
  namespace: secretly
type: Opaque
data:
  url: cG9zdGdyZXNxbDovL3VzZXI6cGFzc3dvcmRAcG9zdGdyZXM6NTQzMi9zZWNyZXRseQ==  # base64 encoded
```

## Monitoreo y Logging

### ServiceMonitor para Prometheus

```yaml
# servicemonitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: secretly-monitor
  namespace: secretly
spec:
  selector:
    matchLabels:
      app: secretly
  endpoints:
  - port: http
    interval: 30s
    path: /api/v1/env
```

### Configuración de Logging

```yaml
# deployment-with-logging.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secretly
  namespace: secretly
spec:
  template:
    spec:
      containers:
      - name: secretly
        image: rodrwan/secretly:latest
        # ... otras configuraciones
        volumeMounts:
        - name: data-volume
          mountPath: /app/data
        - name: logs-volume
          mountPath: /var/log
      volumes:
      - name: data-volume
        persistentVolumeClaim:
          claimName: secretly-pvc
      - name: logs-volume
        emptyDir: {}
```

## Comandos Útiles

```bash
# Verificar estado del despliegue
kubectl get pods -n secretly
kubectl get services -n secretly
kubectl get pvc -n secretly

# Ver logs de un pod
kubectl logs -f deployment/secretly -n secretly

# Ejecutar comandos en un pod
kubectl exec -it deployment/secretly -n secretly -- sh

# Escalar el deployment
kubectl scale deployment secretly --replicas=3 -n secretly

# Verificar eventos
kubectl get events -n secretly --sort-by='.lastTimestamp'

# Port-forward para acceso local
kubectl port-forward service/secretly-service 8080:80 -n secretly

# Hacer backup de los datos
kubectl exec deployment/secretly -n secretly -- tar -czf /tmp/backup.tar.gz /app/data
kubectl cp secretly/secretly-pod:/tmp/backup.tar.gz ./backup.tar.gz
```

## Troubleshooting

### Problema: Pod no inicia

```bash
# Verificar eventos del pod
kubectl describe pod <pod-name> -n secretly

# Verificar logs
kubectl logs <pod-name> -n secretly

# Verificar configuración del deployment
kubectl describe deployment secretly -n secretly
```

### Problema: PVC no se monta

```bash
# Verificar estado del PVC
kubectl describe pvc secretly-pvc -n secretly

# Verificar storage class
kubectl get storageclass

# Verificar eventos del PVC
kubectl get events -n secretly | grep pvc
```

### Problema: Service no funciona

```bash
# Verificar endpoints
kubectl get endpoints -n secretly

# Verificar configuración del service
kubectl describe service secretly-service -n secretly

# Probar conectividad
kubectl run test-pod --image=busybox -it --rm --restart=Never -- wget -O- http://secretly-service
```

## Seguridad

### Network Policies

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: secretly-network-policy
  namespace: secretly
spec:
  podSelector:
    matchLabels:
      app: secretly
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: kube-system
    ports:
    - protocol: TCP
      port: 53
```

### Pod Security Standards

```yaml
# psp.yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: secretly-pdb
  namespace: secretly
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: secretly
```

## Integración con CI/CD

### ArgoCD

```yaml
# argocd-app.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: secretly
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/rodrwan/secretly.git
    targetRevision: HEAD
    path: k8s
  destination:
    server: https://kubernetes.default.svc
    namespace: secretly
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

### Flux

```yaml
# flux-kustomization.yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: secretly
  namespace: flux-system
spec:
  interval: 1m0s
  path: ./k8s
  prune: true
  sourceRef:
    kind: GitRepository
    name: secretly
  targetNamespace: secretly
```

## Helm Chart (Opcional)

Si prefieres usar Helm, puedes crear un chart:

```bash
# Crear estructura del chart
helm create secretly

# Instalar el chart
helm install secretly ./secretly -n secretly

# Actualizar el chart
helm upgrade secretly ./secretly -n secretly

# Desinstalar
helm uninstall secretly -n secretly
``` 