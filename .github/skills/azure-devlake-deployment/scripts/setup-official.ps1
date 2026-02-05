<#
.SYNOPSIS
    Download and set up official Apache DevLake using Docker Compose

.DESCRIPTION
    Downloads docker-compose.yml and env.example from the latest Apache DevLake
    release, generates an encryption secret, and prepares for local deployment.

.PARAMETER TargetDirectory
    Directory where files will be downloaded (default: current directory)

.PARAMETER Version
    Specific version to download (default: latest release)

.EXAMPLE
    .\setup-official.ps1

.EXAMPLE
    .\setup-official.ps1 -TargetDirectory "C:\devlake" -Version "v1.0.2"
#>

param(
    [string]$TargetDirectory = ".",
    [string]$Version = "latest"
)

$ErrorActionPreference = "Stop"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Apache DevLake - Official Setup" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Resolve target directory
$TargetDirectory = Resolve-Path $TargetDirectory -ErrorAction SilentlyContinue
if (-not $TargetDirectory) {
    $TargetDirectory = New-Item -ItemType Directory -Path $TargetDirectory -Force
}
Write-Host "Target directory: $TargetDirectory" -ForegroundColor Yellow

# Get latest release version if not specified
if ($Version -eq "latest") {
    Write-Host "`nFetching latest release version..." -ForegroundColor Yellow
    try {
        $releaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/apache/incubator-devlake/releases/latest" -TimeoutSec 30
        $Version = $releaseInfo.tag_name
        Write-Host "  Latest version: $Version" -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to fetch latest release. Check your internet connection or specify a version manually with -Version."
        exit 1
    }
}

# Construct download URLs
$baseUrl = "https://raw.githubusercontent.com/apache/incubator-devlake/$Version"
$filesToDownload = @(
    @{ Name = "docker-compose.yml"; Url = "$baseUrl/docker-compose.yml" },
    @{ Name = "env.example"; Url = "$baseUrl/env.example" }
)

# Download files
Write-Host "`nDownloading files for $Version..." -ForegroundColor Yellow
foreach ($file in $filesToDownload) {
    $targetPath = Join-Path $TargetDirectory $file.Name
    Write-Host "  Downloading $($file.Name)..." -ForegroundColor Gray
    try {
        Invoke-WebRequest -Uri $file.Url -OutFile $targetPath -TimeoutSec 60
        Write-Host "  ✓ $($file.Name)" -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to download $($file.Name) from $($file.Url)"
        exit 1
    }
}

# Rename env.example to .env
$envExamplePath = Join-Path $TargetDirectory "env.example"
$envPath = Join-Path $TargetDirectory ".env"
if (Test-Path $envPath) {
    Write-Host "`n  .env already exists. Backing up to .env.backup" -ForegroundColor Yellow
    Copy-Item $envPath "$envPath.backup" -Force
}
Move-Item $envExamplePath $envPath -Force
Write-Host "  ✓ Renamed env.example to .env" -ForegroundColor Green

# Generate encryption secret
Write-Host "`nGenerating ENCRYPTION_SECRET..." -ForegroundColor Yellow
$encryptionSecret = -join ((65..90) | Get-Random -Count 128 | ForEach-Object { [char]$_ })

# Update .env file with encryption secret
$envContent = Get-Content $envPath -Raw
if ($envContent -match 'ENCRYPTION_SECRET=') {
    $envContent = $envContent -replace 'ENCRYPTION_SECRET=.*', "ENCRYPTION_SECRET=$encryptionSecret"
} else {
    $envContent += "`nENCRYPTION_SECRET=$encryptionSecret"
}
Set-Content $envPath $envContent -NoNewline
Write-Host "  ✓ ENCRYPTION_SECRET generated and saved" -ForegroundColor Green

# Verify Docker is available
Write-Host "`nChecking Docker..." -ForegroundColor Yellow
try {
    $dockerVersion = docker version --format '{{.Server.Version}}' 2>$null
    Write-Host "  ✓ Docker $dockerVersion found" -ForegroundColor Green
}
catch {
    Write-Host "  ⚠ Docker not found or not running" -ForegroundColor Yellow
    Write-Host "    Please install Docker Desktop: https://docs.docker.com/get-docker" -ForegroundColor Gray
}

# Success message
Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  Setup Complete!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green

Write-Host "Files created in: $TargetDirectory" -ForegroundColor Yellow
Write-Host "  • docker-compose.yml"
Write-Host "  • .env (with ENCRYPTION_SECRET)"
Write-Host ""

Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. cd $TargetDirectory"
Write-Host "  2. docker compose up -d"
Write-Host "  3. Wait 2-3 minutes for services to start"
Write-Host "  4. Open Config UI: http://localhost:4000"
Write-Host "  5. Open Grafana:   http://localhost:3002 (admin/admin)"
Write-Host ""

Write-Host "To stop DevLake:" -ForegroundColor Yellow
Write-Host "  docker compose down"
Write-Host ""

Write-Host "Documentation: https://devlake.apache.org/docs/GettingStarted/DockerComposeSetup" -ForegroundColor Gray
Write-Host ""
