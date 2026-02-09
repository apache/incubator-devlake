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

.PARAMETER RepoUrl
    Clone a remote repository instead of using the local repo.
    Useful for deploying from a fork (e.g., https://github.com/ewega/incubator-devlake)

.PARAMETER UseOfficialImages
    Use official Apache DevLake images from Docker Hub instead of building.
    No ACR needed. Uses apache/devlake:latest, apache/devlake-config-ui:latest, apache/devlake-dashboard:latest

.EXAMPLE
    .\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus"

.EXAMPLE
    .\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus" -SkipImageBuild

.EXAMPLE
    .\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus" -RepoUrl "https://github.com/ewega/incubator-devlake"

.EXAMPLE
    .\deploy.ps1 -ResourceGroupName "devlake-rg" -Location "eastus" -UseOfficialImages
#>

param(
    [Parameter(Mandatory = $true)]
    [string]$ResourceGroupName,

    [Parameter(Mandatory = $true)]
    [string]$Location,

    [string]$BaseName = "devlake",

    [switch]$SkipImageBuild,

    [string]$RepoUrl,

    [switch]$UseOfficialImages
)

$ErrorActionPreference = "Stop"

# Handle UseOfficialImages mode (no repo needed, no build needed)
if ($UseOfficialImages) {
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "  DevLake Azure Deployment (Official Images)" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
    Write-Host "Using official Apache DevLake images from Docker Hub" -ForegroundColor Yellow
    Write-Host "  • apache/devlake:latest" -ForegroundColor Gray
    Write-Host "  • apache/devlake-config-ui:latest" -ForegroundColor Gray
    Write-Host "  • apache/devlake-dashboard:latest" -ForegroundColor Gray
    $RepoRoot = $null
    $SkipImageBuild = $true
} elseif ($RepoUrl) {
    # Clone remote repository
    $CloneDir = Join-Path $env:TEMP "devlake-clone-$(Get-Date -Format 'yyyyMMddHHmmss')"
    Write-Host "Cloning $RepoUrl to $CloneDir..." -ForegroundColor Yellow
    git clone --depth 1 $RepoUrl $CloneDir
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to clone repository: $RepoUrl"
        exit 1
    }
    $RepoRoot = $CloneDir
    Write-Host "  Cloned successfully." -ForegroundColor Green
} else {
    # Find local repo root (look for backend/Dockerfile)
    $RepoRoot = Get-Location
    while ($RepoRoot -and -not (Test-Path "$RepoRoot/backend/Dockerfile")) {
        $RepoRoot = Split-Path $RepoRoot -Parent
    }
    if (-not $RepoRoot -or -not (Test-Path "$RepoRoot/backend/Dockerfile")) {
        Write-Error "Could not find DevLake repository root. Run from within the incubator-devlake directory, use -RepoUrl, or use -UseOfficialImages."
        exit 1
    }
}
if ($RepoRoot) {
    Write-Host "Repo root: $RepoRoot" -ForegroundColor Gray
}

# Generate unique suffix
$UniqueSuffix = (Get-FileHash -InputStream ([IO.MemoryStream]::new([Text.Encoding]::UTF8.GetBytes($ResourceGroupName)))).Hash.Substring(0,5).ToLower()
$AcrName = "devlakeacr$UniqueSuffix"

if (-not $UseOfficialImages) {
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "  DevLake Azure Deployment" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

Write-Host "Configuration:" -ForegroundColor Yellow
Write-Host "  Resource Group: $ResourceGroupName"
Write-Host "  Location: $Location"
Write-Host "  Base Name: $BaseName"
if (-not $UseOfficialImages) {
    Write-Host "  ACR Name: $AcrName"
} else {
    Write-Host "  Images: Official (Docker Hub)"
}
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
    Write-Host "  Using repo root: $RepoRoot"
    Push-Location $RepoRoot
    
    try {
        # First deploy ACR only to get login server
        Write-Host "`nDeploying ACR first..." -ForegroundColor Yellow
        $acrDeployment = az deployment group create `
            --resource-group $ResourceGroupName `
            --template-file ".github/skills/azure-devlake-deployment/bicep/main.bicep" `
            --parameters baseName=$BaseName uniqueSuffix=$UniqueSuffix mysqlAdminPassword=$MysqlPassword encryptionSecret=$EncryptionSecret `
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

# Ensure MySQL is running (Azure auto-stops Burstable tier after creation)
$mysqlName = "${BaseName}mysql${UniqueSuffix}"
Write-Host "`nChecking if MySQL server exists and is running..." -ForegroundColor Yellow
$mysqlState = az mysql flexible-server show --name $mysqlName --resource-group $ResourceGroupName --query state -o tsv 2>$null
if ($mysqlState -eq "Stopped") {
    Write-Host "  MySQL server is stopped (Azure auto-stop). Starting..." -ForegroundColor Yellow
    az mysql flexible-server start --name $mysqlName --resource-group $ResourceGroupName --output none
    Write-Host "  Waiting 30s for MySQL to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 30
    Write-Host "  MySQL started." -ForegroundColor Green
} elseif ($mysqlState) {
    Write-Host "  MySQL state: $mysqlState" -ForegroundColor Green
} else {
    Write-Host "  MySQL not yet created (will be created by Bicep)." -ForegroundColor Gray
}

# Deploy infrastructure
Write-Host "`nDeploying infrastructure with Bicep..." -ForegroundColor Yellow

# Determine template path based on mode
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
if ($UseOfficialImages) {
    $templatePath = Join-Path $scriptDir "main-official.bicep"
    Write-Host "  Using official images template: $templatePath" -ForegroundColor Gray
} else {
    $templatePath = ".github/skills/azure-devlake-deployment/bicep/main.bicep"
    if (-not (Test-Path $templatePath)) {
        $templatePath = Join-Path $RepoRoot $templatePath
    }
    if (-not (Test-Path $templatePath)) {
        $templatePath = Join-Path $scriptDir "main.bicep"
    }
}

if (-not (Test-Path $templatePath)) {
    Write-Error "Could not find Bicep template at: $templatePath"
    exit 1
}

# Deploy with appropriate parameters based on mode
if ($UseOfficialImages) {
    $deployment = az deployment group create `
        --resource-group $ResourceGroupName `
        --template-file $templatePath `
        --parameters `
            baseName=$BaseName `
            uniqueSuffix=$UniqueSuffix `
            mysqlAdminPassword=$MysqlPassword `
            encryptionSecret=$EncryptionSecret `
        --query "properties.outputs" -o json | ConvertFrom-Json
} else {
    $deployment = az deployment group create `
        --resource-group $ResourceGroupName `
        --template-file $templatePath `
        --parameters `
            baseName=$BaseName `
            uniqueSuffix=$UniqueSuffix `
            mysqlAdminPassword=$MysqlPassword `
            encryptionSecret=$EncryptionSecret `
            acrName=$AcrName `
        --query "properties.outputs" -o json | ConvertFrom-Json
}

# Verify deployment outputs
if (-not $deployment) {
    Write-Error "Bicep deployment failed or returned no outputs. Check 'az deployment group show' for details."
    exit 1
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  Deployment Complete!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green

Write-Host "Endpoints:" -ForegroundColor Yellow
Write-Host "  Backend API: $($deployment.backendEndpoint.value)"
Write-Host "  Config UI:   $($deployment.configUiEndpoint.value)"
Write-Host "  Grafana:     $($deployment.grafanaEndpoint.value)"

Write-Host "`nResources:" -ForegroundColor Yellow
if (-not $UseOfficialImages -and $deployment.acrName) {
    Write-Host "  ACR:         $($deployment.acrName.value)"
}
Write-Host "  Key Vault:   $($deployment.keyVaultName.value)"
Write-Host "  MySQL:       $($deployment.mysqlServerName.value)"
if ($UseOfficialImages) {
    Write-Host "  Images:      Official (Docker Hub)"
}

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
if ($RepoRoot) {
    $stateFile = Join-Path $RepoRoot ".devlake-azure.json"
} else {
    $stateFile = Join-Path (Get-Location) ".devlake-azure.json"
}

# Get values from deployment outputs with null-safe access
$backendEndpoint = if ($deployment.backendEndpoint) { $deployment.backendEndpoint.value } else { "http://${BaseName}-${UniqueSuffix}.${Location}.azurecontainer.io:8080" }
$grafanaEndpoint = if ($deployment.grafanaEndpoint) { $deployment.grafanaEndpoint.value } else { "http://${BaseName}-grafana-${UniqueSuffix}.${Location}.azurecontainer.io:3000" }
$configUiEndpoint = if ($deployment.configUiEndpoint) { $deployment.configUiEndpoint.value } else { "http://${BaseName}-ui-${UniqueSuffix}.${Location}.azurecontainer.io:4000" }
$acrNameOutput = if ($deployment.acrName) { $deployment.acrName.value } else { if ($UseOfficialImages) { $null } else { $AcrName } }
$keyVaultName = if ($deployment.keyVaultName) { $deployment.keyVaultName.value } else { "${BaseName}kv${UniqueSuffix}" }
$mysqlServerName = if ($deployment.mysqlServerName) { $deployment.mysqlServerName.value } else { "${BaseName}mysql${UniqueSuffix}" }

$state = @{
    deployedAt = (Get-Date -Format "o")
    method = if ($UseOfficialImages) { "bicep-official" } else { "bicep" }
    subscription = $account.name
    subscriptionId = $account.id
    resourceGroup = $ResourceGroupName
    region = $Location
    suffix = $UniqueSuffix
    useOfficialImages = [bool]$UseOfficialImages
    resources = @{
        acr = $acrNameOutput
        keyVault = $keyVaultName
        mysql = $mysqlServerName
        database = "lake"
        containers = @(
            "$BaseName-backend-$UniqueSuffix",
            "$BaseName-grafana-$UniqueSuffix",
            "$BaseName-ui-$UniqueSuffix"
        )
    }
    endpoints = @{
        backend = $backendEndpoint
        grafana = $grafanaEndpoint
        configUi = $configUiEndpoint
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
