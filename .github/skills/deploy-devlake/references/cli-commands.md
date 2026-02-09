# Azure CLI Commands Reference

## Login & Subscription

```bash
az login
az account list --output table
az account set --subscription "<name-or-id>"
az group list --output table
az group show --name <resource-group>
```

## Resource Group

```bash
az group create --name <resource-group> --location <region>
```

## Container Registry

```bash
# Create
az acr create --name <registry-name> --resource-group <resource-group> --sku Basic --location <region>

# Enable admin access
az acr update --name <registry-name> --admin-enabled true

# Get credentials
az acr credential show --name <registry-name>

# List images
az acr repository list --name <registry-name>

# Login for Docker push
az acr login --name <registry-name>
```

## Key Vault

```bash
# Create with RBAC
az keyvault create \
  --name <keyvault-name> \
  --resource-group <resource-group> \
  --location <region> \
  --enable-rbac-authorization true

# CRITICAL: Assign Secrets Officer role (required before storing secrets)
USER_ID=$(az ad signed-in-user show --query id -o tsv)
KV_ID=$(az keyvault show --name <keyvault-name> --query id -o tsv)
az role assignment create --role "Key Vault Secrets Officer" --assignee $USER_ID --scope $KV_ID

# Store secrets
az keyvault secret set --vault-name <keyvault-name> --name db-admin-password --value "<password>"
az keyvault secret set --vault-name <keyvault-name> --name encryption-secret --value "<32-char-secret>"

# Retrieve secrets
az keyvault secret show --vault-name <keyvault-name> --name db-admin-password --query value -o tsv
```

## MySQL Flexible Server

```bash
# Create server
az mysql flexible-server create \
  --name <mysql-name> \
  --resource-group <resource-group> \
  --location <region> \
  --admin-user merico \
  --admin-password <password> \
  --sku-name Standard_B1ms \
  --tier Burstable \
  --version 8.0.21 \
  --storage-size 32 \
  --public-access 0.0.0.0

# IMPORTANT: Start server (Azure may auto-stop Burstable tier)
az mysql flexible-server start --name <mysql-name> --resource-group <resource-group>

# Verify running
az mysql flexible-server show --name <mysql-name> --resource-group <resource-group> --query state

# Create database
az mysql flexible-server db create \
  --resource-group <resource-group> \
  --server-name <mysql-name> \
  --database-name lake
```

## PostgreSQL Flexible Server

```bash
# Create server
az postgres flexible-server create \
  --name <postgres-name> \
  --resource-group <resource-group> \
  --location <region> \
  --admin-user merico \
  --admin-password <password> \
  --sku-name Standard_B1ms \
  --tier Burstable \
  --version 14 \
  --storage-size 32 \
  --public-access 0.0.0.0

# Create database
az postgres flexible-server db create \
  --resource-group <resource-group> \
  --server-name <postgres-name> \
  --database-name lake
```

## Docker Build & Push

```bash
# Build all images (from repo root)
docker build -t devlake-backend:latest -f backend/Dockerfile ./backend
docker build -t devlake-config-ui:latest -f config-ui/Dockerfile ./config-ui
docker build -t devlake-grafana:latest -f grafana/Dockerfile ./grafana

# Tag for ACR
docker tag devlake-backend:latest <acr>.azurecr.io/devlake-backend:latest
docker tag devlake-config-ui:latest <acr>.azurecr.io/devlake-config-ui:latest
docker tag devlake-grafana:latest <acr>.azurecr.io/devlake-grafana:latest

# Push (set DOCKER_CLIENT_TIMEOUT=600 if timeouts occur)
az acr login --name <acr>
docker push <acr>.azurecr.io/devlake-backend:latest
docker push <acr>.azurecr.io/devlake-config-ui:latest
docker push <acr>.azurecr.io/devlake-grafana:latest
```

## Container Instance Operations

```bash
# Show container details
az container show --name <container-name> --resource-group <rg>

# View logs
az container logs --name <container-name> --resource-group <rg>

# Delete container
az container delete --name <container-name> --resource-group <rg> --yes
```

## AKS Operations

```bash
az aks list --output table
az aks get-credentials --name <cluster> --resource-group <resource-group>
az aks show --name <cluster> --resource-group <resource-group>
az aks scale --name <cluster> --resource-group <resource-group> --node-count 3
az aks check-acr --name <cluster> --acr <registry-name>
```

## Kubernetes Operations

```bash
kubectl get pods -n devlake
kubectl logs -f <pod-name> -n devlake
kubectl describe pod <pod-name> -n devlake
kubectl exec -it <pod-name> -n devlake -- /bin/sh
kubectl get service -n devlake
```
