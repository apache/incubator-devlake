<#
.SYNOPSIS
    Deploy DevLake to Azure using Bicep

.DESCRIPTION
    This script builds Docker images, pushes them to ACR, and deploys 
    all infrastructure using the Bicep template.

.PARAMETER ResourceGroupName
    Name of the Azure Resource Group (will be created if doesn't exist)

.PARAMETER Location
    Azure region (e.g., eastus, westeurope, southafricanorth)

.PARAMETER BaseName
    Base name for resources (default: devlake)

.PARAMETER SkipImageBuild
    Skip Docker image building (use if images already in ACR)

.EXAMPLE
    .\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus"

.EXAMPLE
    .\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus" -SkipImageBuild
#>

param(
    [Parameter(Mandatory = $true)]
    [string]$ResourceGroupName,

    [Parameter(Mandatory = $true)]
    [string]$Location,

    [string]$BaseName = "devlake",

    [switch]$SkipImageBuild
)

$ErrorActionPreference = "Stop"

# Generate unique suffix
$UniqueSuffix = (Get-FileHash -InputStream ([IO.MemoryStream]::new([Text.Encoding]::UTF8.GetBytes($ResourceGroupName)))).Hash.Substring(0,5).ToLower()
$AcrName = "devlakeacr$UniqueSuffix"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  DevLake Azure Deployment" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Configuration:" -ForegroundColor Yellow
Write-Host "  Resource Group: $ResourceGroupName"
Write-Host "  Location: $Location"
Write-Host "  Base Name: $BaseName"
Write-Host "  ACR Name: $AcrName"
Write-Host ""

# Check Azure CLI login
Write-Host "Checking Azure CLI login..." -ForegroundColor Yellow
$account = az account show 2>$null | ConvertFrom-Json
if (-not $account) {
    Write-Host "Not logged in to Azure. Running 'az login'..." -ForegroundColor Yellow
    az login
}
Write-Host "Logged in as: $($account.user.name)" -ForegroundColor Green

# Create Resource Group
Write-Host "`nCreating Resource Group..." -ForegroundColor Yellow
az group create --name $ResourceGroupName --location $Location --output none
Write-Host "Resource Group created." -ForegroundColor Green

# Generate secrets
Write-Host "`nGenerating secrets..." -ForegroundColor Yellow
$MysqlPassword = [guid]::NewGuid().ToString().Replace('-','').Substring(0,16) + "Aa1!"
$EncryptionSecret = [guid]::NewGuid().ToString().Replace('-','').Substring(0,32)
Write-Host "Secrets generated." -ForegroundColor Green

# Build and push images (unless skipped)
if (-not $SkipImageBuild) {
    Write-Host "`nBuilding Docker images..." -ForegroundColor Yellow
    
    # Find repo root (look for backend/Dockerfile)
    $RepoRoot = Get-Location
    while ($RepoRoot -and -not (Test-Path "$RepoRoot/backend/Dockerfile")) {
        $RepoRoot = Split-Path $RepoRoot -Parent
    }
    
    if (-not $RepoRoot -or -not (Test-Path "$RepoRoot/backend/Dockerfile")) {
        Write-Error "Could not find DevLake repository root. Run from within the incubator-devlake directory."
        exit 1
    }
    
    Write-Host "  Found repo root: $RepoRoot"
    Push-Location $RepoRoot
    
    try {
        # First deploy ACR only to get login server
        Write-Host "`nDeploying ACR first..." -ForegroundColor Yellow
        $acrDeployment = az deployment group create `
            --resource-group $ResourceGroupName `
            --template-file ".github/skills/azure-devlake-deployment/bicep/main.bicep" `
            --parameters baseName=$BaseName mysqlAdminPassword=$MysqlPassword encryptionSecret=$EncryptionSecret `
            --query "properties.outputs" -o json 2>$null | ConvertFrom-Json
        
        $AcrLoginServer = $acrDeployment.acrLoginServer.value
        if (-not $AcrLoginServer) {
            # If deployment failed, create ACR manually
            Write-Host "  Creating ACR manually..." -ForegroundColor Yellow
            az acr create --name $AcrName --resource-group $ResourceGroupName --sku Basic --location $Location --admin-enabled true --output none
            $AcrLoginServer = "$AcrName.azurecr.io"
        }
        
        Write-Host "  ACR Login Server: $AcrLoginServer" -ForegroundColor Green
        
        # Login to ACR
        Write-Host "`nLogging into ACR..." -ForegroundColor Yellow
        az acr login --name $AcrName
        
        # Build and push images
        $images = @(
            @{ Name = "devlake-backend"; Dockerfile = "backend/Dockerfile"; Context = "./backend" },
            @{ Name = "devlake-config-ui"; Dockerfile = "config-ui/Dockerfile"; Context = "./config-ui" },
            @{ Name = "devlake-grafana"; Dockerfile = "grafana/Dockerfile"; Context = "./grafana" }
        )
        
        foreach ($img in $images) {
            Write-Host "`n  Building $($img.Name)..." -ForegroundColor Yellow
            docker build -t "$($img.Name):latest" -f $img.Dockerfile $img.Context
            docker tag "$($img.Name):latest" "$AcrLoginServer/$($img.Name):latest"
            
            Write-Host "  Pushing $($img.Name)..." -ForegroundColor Yellow
            $env:DOCKER_CLIENT_TIMEOUT = 600
            docker push "$AcrLoginServer/$($img.Name):latest"
        }
        
        Write-Host "`nAll images pushed successfully." -ForegroundColor Green
    }
    finally {
        Pop-Location
    }
}

# Deploy infrastructure
Write-Host "`nDeploying infrastructure with Bicep..." -ForegroundColor Yellow

$templatePath = ".github/skills/azure-devlake-deployment/bicep/main.bicep"
if (-not (Test-Path $templatePath)) {
    $templatePath = Join-Path $RepoRoot $templatePath
}

$deployment = az deployment group create `
    --resource-group $ResourceGroupName `
    --template-file $templatePath `
    --parameters `
        baseName=$BaseName `
        mysqlAdminPassword=$MysqlPassword `
        encryptionSecret=$EncryptionSecret `
        acrName=$AcrName `
    --query "properties.outputs" -o json | ConvertFrom-Json

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  Deployment Complete!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green

Write-Host "Endpoints:" -ForegroundColor Yellow
Write-Host "  Backend API: $($deployment.backendEndpoint.value)"
Write-Host "  Config UI:   $($deployment.configUiEndpoint.value)"
Write-Host "  Grafana:     $($deployment.grafanaEndpoint.value)"

Write-Host "`nResources:" -ForegroundColor Yellow
Write-Host "  ACR:         $($deployment.acrName.value)"
Write-Host "  Key Vault:   $($deployment.keyVaultName.value)"
Write-Host "  MySQL:       $($deployment.mysqlServerName.value)"

Write-Host "`nSecrets (save these!):" -ForegroundColor Red
Write-Host "  MySQL Password:     $MysqlPassword"
Write-Host "  Encryption Secret:  $EncryptionSecret"

# Wait for backend and trigger migration
Write-Host "`nWaiting for backend to start..." -ForegroundColor Yellow
$backendUrl = $deployment.backendEndpoint.value
$maxAttempts = 30
$attempt = 0
$backendReady = $false

while ($attempt -lt $maxAttempts -and -not $backendReady) {
    $attempt++
    try {
        $response = Invoke-WebRequest -Uri "$backendUrl/ping" -TimeoutSec 5 -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            $backendReady = $true
            Write-Host "  Backend is responding!" -ForegroundColor Green
        }
    }
    catch {
        Write-Host "  Attempt $attempt/$maxAttempts - waiting..." -ForegroundColor Gray
        Start-Sleep -Seconds 10
    }
}

if ($backendReady) {
    Write-Host "`nTriggering database migration..." -ForegroundColor Yellow
    try {
        $migrationResponse = Invoke-RestMethod -Uri "$backendUrl/proceed-db-migration" -Method GET -TimeoutSec 120
        Write-Host "  Migration completed successfully!" -ForegroundColor Green
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -eq 428) {
            Write-Host "  Migration confirmation required. Triggering..." -ForegroundColor Yellow
            $migrationResponse = Invoke-RestMethod -Uri "$backendUrl/proceed-db-migration" -Method GET -TimeoutSec 120
            Write-Host "  Migration completed!" -ForegroundColor Green
        }
        elseif ($statusCode -eq 200 -or $null -eq $statusCode) {
            Write-Host "  Migration already complete or not needed." -ForegroundColor Green
        }
        else {
            Write-Host "  Migration may need manual triggering. Status: $statusCode" -ForegroundColor Yellow
            Write-Host "  Run: Invoke-RestMethod -Uri '$backendUrl/proceed-db-migration' -Method GET" -ForegroundColor Yellow
        }
    }
}
else {
    Write-Host "  Backend not ready after $maxAttempts attempts." -ForegroundColor Yellow
    Write-Host "  You may need to manually trigger migration once backend starts:" -ForegroundColor Yellow
    Write-Host "  Invoke-RestMethod -Uri '$backendUrl/proceed-db-migration' -Method GET" -ForegroundColor Yellow
}

# Save deployment state file
$stateFile = Join-Path $RepoRoot ".devlake-azure.json"
$state = @{
    deployedAt = (Get-Date -Format "o")
    method = "bicep"
    subscription = $account.name
    subscriptionId = $account.id
    resourceGroup = $ResourceGroupName
    region = $Location
    suffix = $UniqueSuffix
    resources = @{
        acr = $deployment.acrName.value
        keyVault = $deployment.keyVaultName.value
        mysql = $deployment.mysqlServerName.value
        database = "lake"
        containers = @(
            "$BaseName-backend-$UniqueSuffix",
            "$BaseName-grafana-$UniqueSuffix",
            "$BaseName-ui-$UniqueSuffix"
        )
    }
    endpoints = @{
        backend = $deployment.backendEndpoint.value
        grafana = $deployment.grafanaEndpoint.value
        configUi = $deployment.configUiEndpoint.value
    }
    secrets = @{
        keyVaultSecrets = @("db-admin-password", "encryption-secret")
        mysqlPassword = $MysqlPassword
        encryptionSecret = $EncryptionSecret
    }
}
$state | ConvertTo-Json -Depth 5 | Set-Content $stateFile -Encoding UTF8
Write-Host "`nState saved to: $stateFile" -ForegroundColor Green

Write-Host "`nNext Steps:" -ForegroundColor Cyan
Write-Host "  1. Wait 2-3 minutes for containers to start"
Write-Host "  2. Open Config UI: $($deployment.configUiEndpoint.value)"
Write-Host "  3. Configure your data sources"
Write-Host ""
Write-Host "To cleanup later, run:" -ForegroundColor Yellow
Write-Host "  az group delete --name $ResourceGroupName --yes --no-wait"
Write-Host ""
