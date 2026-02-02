# Troubleshooting Guide

## Container Fails to Start

**Symptoms:** Container in "Waiting" state, restarts repeatedly

**Debug:**
```bash
# View logs
az container logs --name <container> --resource-group <rg>

# Create with no restart to see error
az container create ... --restart-policy Never
```

**Common causes:**
- Missing environment variables
- Image pull failure
- Database connection failure

---

## Backend Panics Immediately

**Error:** `panic: runtime error` or container exits immediately

**Cause:** Missing `ENCRYPTION_SECRET`

**Fix:** Set `ENCRYPTION_SECRET` environment variable with exactly 32 characters.

**Generate:**
```bash
# Bash
openssl rand -base64 24 | tr -dc 'a-zA-Z0-9' | head -c 32

# PowerShell
[guid]::NewGuid().ToString().Replace('-','').Substring(0,32)
```

---

## Datetime Scan Error

**Error:** 
```
unsupported Scan, storing driver.Value type []uint8 into type *time.Time
```

**Cause:** DB_URL missing `parseTime=True&loc=UTC`

**Fix:** Add these query parameters to your DB_URL:
```
mysql://user:pass@host:3306/db?charset=utf8mb4&parseTime=True&loc=UTC&tls=true
```

---

## Key Vault Forbidden Error

**Error:** `Forbidden` when running `az keyvault secret set`

**Cause:** Missing RBAC role assignment

**Fix:**
```bash
USER_ID=$(az ad signed-in-user show --query id -o tsv)
KV_ID=$(az keyvault show --name <kv-name> --query id -o tsv)
az role assignment create --role "Key Vault Secrets Officer" --assignee $USER_ID --scope $KV_ID
```

Wait ~30 seconds for propagation, then retry.

---

## MySQL Server Stopped

**Error:** Connection refused to MySQL

**Cause:** Azure auto-stops Burstable tier servers after inactivity

**Fix:**
```bash
# Start server
az mysql flexible-server start --name <mysql-name> --resource-group <rg>

# Verify running
az mysql flexible-server show --name <mysql-name> --resource-group <rg> --query state
```

---

## Database Connection Fails

**Symptoms:** Backend can't connect to database

**Checks:**
1. **Server running?**
   ```bash
   az mysql flexible-server show --name <name> --resource-group <rg> --query state
   ```

2. **Firewall allows access?**
   - For testing: `--public-access 0.0.0.0` allows all IPs
   - For production: Add specific ACI/AKS IPs

3. **TLS enabled in connection string?**
   - Must include `tls=true` for Azure MySQL

4. **Correct credentials?**
   - Test with mysql CLI from local machine

---

## Cannot Pull Image from ACR

**Symptoms:** Container stuck in "Waiting", image pull errors

**For ACI:**
```bash
# Verify ACR credentials
az acr credential show --name <acr>

# Verify image exists
az acr repository list --name <acr>
```

**For AKS:**
```bash
# Check ACR integration
az aks check-acr --name <cluster> --acr <acr>

# Re-attach if needed
az aks update --name <cluster> --resource-group <rg> --attach-acr <acr>
```

---

## External IP Pending

**Symptoms:** `kubectl get service` shows `<pending>` for EXTERNAL-IP

**Causes:**
1. Azure still provisioning (wait 2-5 minutes)
2. Region quota exceeded
3. Service type not LoadBalancer

**Check quotas:**
```bash
az network lb list --resource-group MC_<rg>_<cluster>_<region>
```

---

## Docker Push Timeout

**Error:** Timeout during `docker push`

**Fix:** Set timeout:
```bash
export DOCKER_CLIENT_TIMEOUT=600
docker push <image>
```

Or push layers individually by pushing smaller images first.

---

## MySQL Migration Fails with Primary Key Error

**Error:**
```
Error 3750: Unable to create or change a table without a primary key, 
when the system variable 'sql_require_primary_key' is set...
```

**Or:**
```
Error 4108: Cannot drop or rename a column belonging to an invisible primary key
```

**Cause:** Azure MySQL Flexible Server has `sql_generate_invisible_primary_key=ON` by default. DevLake migrations that drop/recreate primary keys fail.

**Fix:**
```bash
az mysql flexible-server parameter set \
  --resource-group <rg> \
  --server-name <mysql-name> \
  --name sql_generate_invisible_primary_key \
  --value OFF

# Restart backend container after changing
az container restart --name devlake-backend --resource-group <rg>
```

---

## Database Migration Required (HTTP 428)

**Symptoms:** Backend returns 428 status, Config UI shows "Database migration required"

**Cause:** Fresh deployment detected schema changes needing migration

**Fix:**
```powershell
# Trigger migration via API
Invoke-RestMethod -Uri "http://<backend-fqdn>:8080/proceed-db-migration" -Method GET
```

Wait for migration to complete (check logs), then restart Config UI if needed.

---

## Config UI Returns 500 Error

**Symptoms:** Config UI accessible but shows 500 Internal Server Error

**Cause:** nginx in Config UI adds `http://` to all proxy upstreams. If `DEVLAKE_ENDPOINT` or `GRAFANA_ENDPOINT` already contains `http://`, the result is invalid (`http://http://...`).

**Fix:** Ensure environment variables do NOT include the protocol prefix:
```
DEVLAKE_ENDPOINT=hostname:8080     ✓ Correct
DEVLAKE_ENDPOINT=http://hostname:8080  ✗ Wrong - nginx will make this http://http://...
```

**To fix existing deployment:**
```bash
# Delete and recreate Config UI container with correct endpoint format
az container delete --name devlake-config-ui --resource-group <rg> --yes

az container create \
  --resource-group <rg> \
  --name devlake-config-ui \
  --image apache/devlake-config-ui:v1.0.2 \
  --ports 4000 \
  --ip-address Public \
  --environment-variables \
    DEVLAKE_ENDPOINT=<backend-fqdn>:8080 \
    GRAFANA_ENDPOINT=<grafana-fqdn>:3000
```
