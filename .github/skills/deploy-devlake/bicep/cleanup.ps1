<#
.SYNOPSIS
    Cleanup DevLake Azure deployment

.DESCRIPTION
    Deletes all Azure resources created by deploy.ps1.
    Reads the state file to determine what to delete.

.PARAMETER Force
    Skip confirmation prompt

.PARAMETER StateFile
    Path to state file (default: .devlake-azure.json in repo root)

.PARAMETER KeepResourceGroup
    Delete resources but keep the resource group

.EXAMPLE
    .\cleanup.ps1

.EXAMPLE
    .\cleanup.ps1 -Force
#>

param(
    [switch]$Force,
    [string]$StateFile,
    [switch]$KeepResourceGroup
)

$ErrorActionPreference = "Stop"

# Find state file
if (-not $StateFile) {
    $RepoRoot = Get-Location
    while ($RepoRoot -and -not (Test-Path "$RepoRoot/.devlake-azure.json")) {
        $RepoRoot = Split-Path $RepoRoot -Parent
    }
    
    if ($RepoRoot) {
        $StateFile = Join-Path $RepoRoot ".devlake-azure.json"
    }
}

if (-not (Test-Path $StateFile)) {
    Write-Error "State file not found. Cannot determine what to cleanup.`nExpected: .devlake-azure.json in repo root`nOr specify with: -StateFile <path>"
    exit 1
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  DevLake Azure Cleanup" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Load state
$state = Get-Content $StateFile | ConvertFrom-Json

Write-Host "Deployment found:" -ForegroundColor Yellow
Write-Host "  Deployed: $($state.deployedAt)"
Write-Host "  Resource Group: $($state.resourceGroup)"
Write-Host "  Region: $($state.region)"
Write-Host ""

Write-Host "Resources to delete:" -ForegroundColor Yellow
Write-Host "  Container Registry: $($state.resources.acr)"
Write-Host "  Key Vault: $($state.resources.keyVault)"
Write-Host "  MySQL Server: $($state.resources.mysql)"
Write-Host "  Containers: $($state.resources.containers -join ', ')"
Write-Host ""

Write-Host "Endpoints that will be removed:" -ForegroundColor Yellow
Write-Host "  Backend: $($state.endpoints.backend)"
Write-Host "  Config UI: $($state.endpoints.configUi)"
Write-Host "  Grafana: $($state.endpoints.grafana)"
Write-Host ""

if (-not $Force) {
    $confirm = Read-Host "Are you sure you want to delete ALL these resources? (yes/no)"
    if ($confirm -ne "yes") {
        Write-Host "Cleanup cancelled." -ForegroundColor Yellow
        exit 0
    }
}

# Check Azure login
Write-Host "Checking Azure CLI login..." -ForegroundColor Yellow
$account = az account show 2>$null | ConvertFrom-Json
if (-not $account) {
    Write-Host "Not logged in to Azure. Running 'az login'..." -ForegroundColor Yellow
    az login
}

if ($KeepResourceGroup) {
    # Delete individual resources
    Write-Host "`nDeleting containers..." -ForegroundColor Yellow
    foreach ($container in $state.resources.containers) {
        Write-Host "  Deleting $container..."
        az container delete --name $container --resource-group $state.resourceGroup --yes 2>$null
    }
    
    Write-Host "`nDeleting MySQL server..." -ForegroundColor Yellow
    az mysql flexible-server delete --name $state.resources.mysql --resource-group $state.resourceGroup --yes 2>$null
    
    Write-Host "`nDeleting Container Registry..." -ForegroundColor Yellow
    az acr delete --name $state.resources.acr --resource-group $state.resourceGroup --yes 2>$null
    
    Write-Host "`nDeleting Key Vault (soft delete)..." -ForegroundColor Yellow
    az keyvault delete --name $state.resources.keyVault --resource-group $state.resourceGroup 2>$null
    
    Write-Host "`nResource group '$($state.resourceGroup)' kept." -ForegroundColor Green
} else {
    # Delete entire resource group
    Write-Host "`nDeleting resource group '$($state.resourceGroup)'..." -ForegroundColor Yellow
    Write-Host "  This may take a few minutes..." -ForegroundColor Gray
    az group delete --name $state.resourceGroup --yes --no-wait
    Write-Host "  Deletion initiated (running in background)." -ForegroundColor Green
}

# Remove state file
Write-Host "`nRemoving state file..." -ForegroundColor Yellow
Remove-Item $StateFile -Force
Write-Host "State file removed." -ForegroundColor Green

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  Cleanup Complete!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green

if (-not $KeepResourceGroup) {
    Write-Host "Note: Resource group deletion runs in background." -ForegroundColor Yellow
    Write-Host "Check status with: az group show --name $($state.resourceGroup) 2>`$null || echo 'Deleted'" -ForegroundColor Gray
}
Write-Host ""
