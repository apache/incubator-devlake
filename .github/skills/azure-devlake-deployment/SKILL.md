---
name: azure-devlake-deployment
description: Deploy Apache DevLake to Azure using ACI or AKS. Use when deploying DevLake to Azure, creating Azure resources, building Docker images, or troubleshooting Azure DevLake deployments.
---

# Azure DevLake Deployment Skill

Deploy Apache DevLake to Microsoft Azure with either Infrastructure-as-Code (Bicep) or step-by-step CLI commands.

## Quick Start

### Option A: One-Command Deployment (Bicep)

```powershell
# From the incubator-devlake repo root
.\.github\skills\azure-devlake-deployment\bicep\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus"
```

This builds images, creates all Azure resources, and deploys 3 containers. Takes ~10-15 minutes.

### Option B: Step-by-Step CLI

Follow this order:
1. Create Resource Group → [cli-commands.md](references/cli-commands.md)
2. Create ACR, Key Vault, MySQL → [cli-commands.md](references/cli-commands.md)
3. Build & push Docker images → [cli-commands.md](references/cli-commands.md)
4. Deploy containers → [aci-deployment.md](references/aci-deployment.md) or [aks-deployment.md](references/aks-deployment.md)
5. Verify endpoints return 200

### Option C: Using the Custom Agent

For guided interactive deployment with confirmation prompts:
```
@DevlakeAzureDeployer deploy DevLake to Azure
```

## Prerequisites

- Azure CLI installed and logged in (`az account show`)
- Docker installed locally
- Active Azure subscription with contributor permissions

## Deployment Methods

### Option 1: Bicep (Recommended)

Single-command Infrastructure-as-Code deployment.

**Advantages:**
- One command deploys all resources
- No PowerShell escaping issues with DB_URL
- Repeatable, idempotent, version-controlled
- Easier to modify and redeploy

**Files:**
- [bicep/main.bicep](bicep/main.bicep) - Infrastructure template
- [bicep/deploy.ps1](bicep/deploy.ps1) - Deployment script

### Option 2: CLI Commands

Step-by-step Azure CLI commands for manual control.

**Advantages:**
- Full visibility at each step
- Easier to debug individual resources
- Good for learning Azure

**Files:** See [references/](references/) folder

## Prerequisites

- Azure CLI installed and logged in (`az account show`)
- Docker installed locally
- Active Azure subscription with contributor permissions

## Critical Requirements

These apply to **ALL** deployment methods:

| Requirement | Why It Matters |
|-------------|----------------|
| `parseTime=True&loc=UTC&tls=true` in DB_URL | Without this, datetime fields fail to scan |
| `ENCRYPTION_SECRET` (32 chars) | Backend panics immediately without it |
| Key Vault RBAC role | "Forbidden" error when storing secrets |
| Start MySQL after creation | Azure auto-stops Burstable tier servers |

See [references/environment-variables.md](references/environment-variables.md) for complete details.

## Reference Files

| File | Content |
|------|---------|
| [cli-commands.md](references/cli-commands.md) | Resource creation commands |
| [aci-deployment.md](references/aci-deployment.md) | Container instance deployment |
| [aks-deployment.md](references/aks-deployment.md) | Kubernetes deployment |
| [environment-variables.md](references/environment-variables.md) | Required env vars & formats |
| [troubleshooting.md](references/troubleshooting.md) | Common issues & fixes |
| [powershell-escaping.md](references/powershell-escaping.md) | Windows-specific escaping |
| [cleanup.md](references/cleanup.md) | Teardown and resource deletion |
| [k8s-manifest.yaml](references/k8s-manifest.yaml) | K8s deployment template |

## State File

Deployments create `.devlake-azure.json` in the repo root containing:
- Resource names and IDs
- Endpoints
- Secrets (for reference)

Use this file to track what was deployed and for cleanup. Add to `.gitignore`.

### Quick Cleanup

```powershell
# One-command cleanup using state file
.\.github\skills\azure-devlake-deployment\bicep\cleanup.ps1
```

Or manually:
```powershell
$state = Get-Content .devlake-azure.json | ConvertFrom-Json
az group delete --name $state.resourceGroup --yes --no-wait
Remove-Item .devlake-azure.json
```

See [cleanup.md](references/cleanup.md) for selective deletion options.

## Cost Estimate

~$50-75/month (MySQL B1ms + 3 containers + ACR Basic + Key Vault)

## Deployment Plan Template

Present before creating resources:

```
═══════════════════════════════════════════════════════════════
                    DEPLOYMENT PLAN SUMMARY
═══════════════════════════════════════════════════════════════
📍 Target: <subscription> / <region>
🔧 Method: Bicep / CLI

📦 Resources: RG, ACR, Key Vault, MySQL, 3x Container Instances
💰 Cost: ~$50-75/month

🌐 Endpoints:
   • Backend: http://devlake-<suffix>.<region>.azurecontainer.io:8080
   • Config UI: http://devlake-ui-<suffix>.<region>.azurecontainer.io:4000
   • Grafana: http://devlake-grafana-<suffix>.<region>.azurecontainer.io:3000
═══════════════════════════════════════════════════════════════
```
