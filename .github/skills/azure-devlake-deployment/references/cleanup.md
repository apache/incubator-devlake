# Cleanup and Teardown

## Quick Cleanup (Delete Everything)

Delete the entire resource group to remove all deployed resources:

```bash
az group delete --name <resource-group> --yes --no-wait
```

This removes:
- Container Instances (backend, grafana, config-ui)
- MySQL Flexible Server
- Key Vault
- Container Registry
- All associated resources

## Selective Cleanup

### Delete Container Instances Only

```bash
az container delete --name devlake-backend-<suffix> --resource-group <rg> --yes
az container delete --name devlake-grafana-<suffix> --resource-group <rg> --yes
az container delete --name devlake-ui-<suffix> --resource-group <rg> --yes
```

### Delete MySQL Server

```bash
az mysql flexible-server delete --name <mysql-name> --resource-group <rg> --yes
```

### Delete Container Registry

```bash
az acr delete --name <acr-name> --resource-group <rg> --yes
```

### Delete Key Vault

```bash
# Soft delete (recoverable for 90 days)
az keyvault delete --name <kv-name> --resource-group <rg>

# Purge permanently (irreversible)
az keyvault purge --name <kv-name> --location <region>
```

## Using the State File

If you deployed using `deploy.ps1` or the custom agent, a state file `.devlake-azure.json` was created in the repo root.

### Read State File

```powershell
$state = Get-Content .devlake-azure.json | ConvertFrom-Json
$state.resourceGroup  # Get resource group name
$state.resources.containers  # List container names
$state.endpoints.configUi  # Get Config UI URL
```

### Cleanup Using State File

```powershell
$state = Get-Content .devlake-azure.json | ConvertFrom-Json
az group delete --name $state.resourceGroup --yes --no-wait
Remove-Item .devlake-azure.json
```

## AKS Cleanup

### Delete AKS Cluster

```bash
az aks delete --name devlake-aks --resource-group <rg> --yes --no-wait
```

### Delete Kubernetes Resources Only

```bash
kubectl delete namespace devlake
```

## Verify Cleanup

```bash
# List remaining resources in group
az resource list --resource-group <rg> --output table

# Check if resource group still exists
az group show --name <rg> 2>/dev/null || echo "Resource group deleted"
```

## Cost Implications

- **Container Instances**: Billed per-second, stops immediately on delete
- **MySQL**: Billed hourly, stops on delete (may have backup retention costs)
- **ACR**: Billed for storage, stops on delete
- **Key Vault**: Soft-deleted vaults don't incur charges, but purging is irreversible

## Troubleshooting Cleanup

### Resource Group Delete Stuck

```bash
# Check for locks
az lock list --resource-group <rg>

# Remove lock if exists
az lock delete --name <lock-name> --resource-group <rg>
```

### Key Vault Soft Delete

Azure Key Vault uses soft-delete by default. To fully remove:

```bash
# List soft-deleted vaults
az keyvault list-deleted

# Purge specific vault
az keyvault purge --name <kv-name> --location <region>
```

### Orphaned Resources

If resource group delete times out, some resources may be orphaned:

```bash
# Force delete all resources in group
az resource list --resource-group <rg> --query "[].id" -o tsv | xargs -L1 az resource delete --ids
```
