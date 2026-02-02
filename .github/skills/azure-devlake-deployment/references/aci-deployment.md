# Azure Container Instances Deployment

Deploy DevLake as three separate containers on Azure Container Instances.

## Prerequisites

- ACR created with images pushed
- Key Vault with secrets stored
- MySQL/PostgreSQL server running with database created

## Get ACR Credentials

```bash
ACR_USER=$(az acr credential show --name <acr> --query username -o tsv)
ACR_PASS=$(az acr credential show --name <acr> --query "passwords[0].value" -o tsv)
```

## Get Secrets from Key Vault

```bash
DB_PASSWORD=$(az keyvault secret show --vault-name <kv> --name db-admin-password --query value -o tsv)
ENCRYPTION_SECRET=$(az keyvault secret show --vault-name <kv> --name encryption-secret --query value -o tsv)
```

## Deploy Backend Container

```bash
az container create \
  --name devlake-backend-<suffix> \
  --resource-group <rg> \
  --image <acr>.azurecr.io/devlake-backend:latest \
  --registry-login-server <acr>.azurecr.io \
  --registry-username $ACR_USER \
  --registry-password $ACR_PASS \
  --dns-name-label devlake-<suffix> \
  --ports 8080 \
  --cpu 2 \
  --memory 4 \
  --environment-variables \
    DB_URL="mysql://merico:${DB_PASSWORD}@<mysql>.mysql.database.azure.com:3306/lake?charset=utf8mb4&parseTime=True&loc=UTC&tls=true" \
    ENCRYPTION_SECRET="$ENCRYPTION_SECRET" \
    PORT=8080 \
    MODE=release \
    PLUGIN_DIR=bin/plugins \
    REMOTE_PLUGIN_DIR=python/plugins \
    LOGGING_DIR=/app/logs \
    TZ=UTC
```

**Endpoint:** `http://devlake-<suffix>.<region>.azurecontainer.io:8080`

## Deploy Grafana Container

```bash
az container create \
  --name devlake-grafana-<suffix> \
  --resource-group <rg> \
  --image <acr>.azurecr.io/devlake-grafana:latest \
  --registry-login-server <acr>.azurecr.io \
  --registry-username $ACR_USER \
  --registry-password $ACR_PASS \
  --dns-name-label devlake-grafana-<suffix> \
  --ports 3000 \
  --cpu 1 \
  --memory 2 \
  --environment-variables \
    GF_SERVER_ROOT_URL="http://devlake-grafana-<suffix>.<region>.azurecontainer.io:3000"
```

**Endpoint:** `http://devlake-grafana-<suffix>.<region>.azurecontainer.io:3000`

## Deploy Config UI Container

```bash
az container create \
  --name devlake-ui-<suffix> \
  --resource-group <rg> \
  --image <acr>.azurecr.io/devlake-config-ui:latest \
  --registry-login-server <acr>.azurecr.io \
  --registry-username $ACR_USER \
  --registry-password $ACR_PASS \
  --dns-name-label devlake-ui-<suffix> \
  --ports 4000 \
  --cpu 1 \
  --memory 2 \
  --environment-variables \
    DEVLAKE_ENDPOINT="http://devlake-<suffix>.<region>.azurecontainer.io:8080" \
    GRAFANA_ENDPOINT="http://devlake-grafana-<suffix>.<region>.azurecontainer.io:3000"
```

**Endpoint:** `http://devlake-ui-<suffix>.<region>.azurecontainer.io:4000`

## Verify Deployment

```bash
# Check container status
az container show --name devlake-backend-<suffix> --resource-group <rg> --query instanceView.state

# Test endpoints
curl http://devlake-<suffix>.<region>.azurecontainer.io:8080/ping
curl http://devlake-grafana-<suffix>.<region>.azurecontainer.io:3000/api/health
curl http://devlake-ui-<suffix>.<region>.azurecontainer.io:4000
```

## Debugging

```bash
# View logs
az container logs --name devlake-backend-<suffix> --resource-group <rg>

# Create with restart-policy Never to see startup errors
az container create ... --restart-policy Never
```
