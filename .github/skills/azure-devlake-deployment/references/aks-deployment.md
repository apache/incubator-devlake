# Azure Kubernetes Service Deployment

Deploy DevLake on AKS for production-grade scalability.

## Create AKS Cluster

```bash
az aks create \
  --name devlake-aks \
  --resource-group <rg> \
  --location <region> \
  --node-count 2 \
  --node-vm-size Standard_D2s_v3 \
  --enable-managed-identity \
  --attach-acr <acr> \
  --generate-ssh-keys

# Get credentials
az aks get-credentials --name devlake-aks --resource-group <rg>
```

## Deploy with Manifest

Use the template at [k8s-manifest.yaml](k8s-manifest.yaml).

Replace placeholders:
- `<password>` - Database password
- `<server>` - MySQL server name
- `<acr>` - ACR name
- `<32-character-random-string>` - ENCRYPTION_SECRET

```bash
kubectl apply -f k8s-manifest.yaml

# Wait for pods
kubectl wait --for=condition=ready pod -l app=devlake-server -n devlake --timeout=300s

# Get external IPs
kubectl get service -n devlake
```

## Verify Deployment

```bash
# Check pods
kubectl get pods -n devlake

# View logs
kubectl logs -f deployment/devlake-server -n devlake

# Describe pod for events
kubectl describe pod -l app=devlake-server -n devlake
```

## Configure Ingress (Optional)

Install NGINX Ingress Controller:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
```

Create ingress resource:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: devlake-ingress
  namespace: devlake
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: devlake.yourdomain.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: devlake-service
            port:
              number: 80
      - path: /
        pathType: Prefix
        backend:
          service:
            name: devlake-ui-service
            port:
              number: 80
```

## Scaling

```bash
# Scale deployment
kubectl scale deployment devlake-server -n devlake --replicas=3

# Scale AKS nodes
az aks scale --name devlake-aks --resource-group <rg> --node-count 3
```

## Troubleshooting

```bash
# Exec into pod
kubectl exec -it <pod-name> -n devlake -- /bin/sh

# Check ACR connectivity
az aks check-acr --name devlake-aks --acr <acr>

# Restart deployment
kubectl rollout restart deployment devlake-server -n devlake
```
