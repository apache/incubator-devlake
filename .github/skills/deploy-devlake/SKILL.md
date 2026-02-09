---
name: azure-devlake-deployment
description: Deploy Apache DevLake to Azure using ACI or AKS. Use when deploying DevLake to Azure, creating Azure resources, building Docker images, or troubleshooting Azure DevLake deployments.
---

# Azure DevLake Deployment Skill

Deploy Apache DevLake with two paths: **Official** (local or Azure) or **Custom** (build from source → Azure).

## Deployment Paths

```
═══════════════════════════════════════════════════════════════
              DEVLAKE DEPLOYMENT - SELECT PATH
═══════════════════════════════════════════════════════════════
1️⃣  Official Apache DevLake (latest release, official images)
    a) Local Docker - quick setup on your machine
    b) Deploy to Azure - ACI containers with managed MySQL

2️⃣  Custom DevLake (build from source → Azure)
    a) Clone a remote fork (e.g., ewega/incubator-devlake)
    b) Use a local repository path
    Build dev images → Deploy to Azure ACI/AKS

☁️  Cloud support: Azure only (AWS/GCP coming soon)
═══════════════════════════════════════════════════════════════
```

## Path 1a: Official Apache DevLake (Local Docker)

Quick local setup using official release images. No Azure resources needed.

```powershell
# Download and set up official DevLake
.\.github\skills\azure-devlake-deployment\scripts\setup-official.ps1

# Or specify a target directory and version
.\.github\skills\azure-devlake-deployment\scripts\setup-official.ps1 -TargetDirectory "C:\devlake" -Version "v1.0.2"

# Then start DevLake
cd C:\devlake
docker compose up -d
```

**Endpoints after startup:**
- Config UI: http://localhost:4000
- Grafana: http://localhost:3002 (admin/admin)
- Backend API: http://localhost:8080

## Path 1b: Official Apache DevLake (Azure Deployment)

Deploy official release images to Azure. No build required, no ACR needed.

```powershell
# Deploy official images to Azure
.\.github\skills\azure-devlake-deployment\bicep\deploy.ps1 `
    -ResourceGroupName "devlake-rg" `
    -Location "eastus" `
    -UseOfficialImages
```

**Cost:** ~$30-50/month (no ACR)

## Path 2: Custom DevLake (Azure Deployment)

Build from source (local repo or remote fork) and deploy to Azure.

### Option A: Deploy from Local Repository (Bicep)

```powershell
# From the incubator-devlake repo root
.\.github\skills\azure-devlake-deployment\bicep\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus"
```

This builds images, creates all Azure resources, and deploys 3 containers. Takes ~10-15 minutes.

### Option B: Deploy from Remote Fork

```powershell
# Clone and deploy a fork
.\.github\skills\azure-devlake-deployment\bicep\deploy.ps1 `
    -ResourceGroupName "devlake-rg" `
    -Location "eastus" `
    -RepoUrl "https://github.com/ewega/incubator-devlake"
```

### Option C: Step-by-Step CLI

Follow this order:
1. Create Resource Group → [cli-commands.md](references/cli-commands.md)
2. Create ACR, Key Vault, MySQL → [cli-commands.md](references/cli-commands.md)
3. Build & push Docker images → [cli-commands.md](references/cli-commands.md)
4. Deploy containers → [aci-deployment.md](references/aci-deployment.md) or [aks-deployment.md](references/aks-deployment.md)
5. Verify endpoints return 200

### Option D: Using the Custom Agent

For guided interactive deployment with confirmation prompts:
```
@DeployDevlake deploy DevLake to Azure
```

**Cost:** ~$50-75/month (includes ACR)

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

| Mode | Monthly Cost | Includes |
|------|-------------|----------|
| Official (Azure) | ~$30-50 | MySQL B1ms + 3 containers + Key Vault |
| Custom (Azure) | ~$50-75 | MySQL B1ms + 3 containers + ACR Basic + Key Vault |

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
