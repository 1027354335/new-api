# New API Dev/Build Helper Script (dev.ps1)
# Usage: .\dev.ps1 [setup|build|build-windows|run|all]

param(
    [Parameter(Position=0)]
    [ValidateSet("setup", "build", "build-windows", "run", "all")]
    [string]$Action = "build"
)

$ErrorActionPreference = "Stop"

# Get version and handle potential encoding/null issues
$VERSION = "unknown"
if (Test-Path VERSION) {
    $rawVersion = Get-Content VERSION -Raw
    if ($null -ne $rawVersion) {
        $VERSION = $rawVersion.Trim()
        if ($VERSION -eq "") { $VERSION = "0.0.1-dev" }
    }
}

function Write-P8-Header($msg) {
    if ($null -eq $msg) { $msg = "NOTIFICATION" }
    Write-Host "`n+----------------------------------------------------------------+" -ForegroundColor Cyan
    Write-Host ("| " + $msg.PadRight(62) + " |") -ForegroundColor Cyan
    Write-Host "+----------------------------------------------------------------+" -ForegroundColor Cyan
}

switch ($Action) {
    "setup" {
        Write-P8-Header "INITIATING SETUP: ALIGNING RESOURCES"
        Write-Host ">> Downloading Go modules..."
        go mod download
        
        Write-Host ">> Installing frontend dependencies (web/default)..."
        Push-Location web/default
        bun install
        Pop-Location
        
        Write-Host ">> Installing frontend dependencies (web/classic)..."
        Push-Location web/classic
        bun install
        Pop-Location
        
        if (-not (Test-Path .env)) {
            Write-Host ">> Creating .env from example..." -ForegroundColor Yellow
            Copy-Item .env.example .env
            Write-Host "!! IMPORTANT: Please configure SQL_DSN in .env before running." -ForegroundColor Red
        }
        Write-Host "`n[SUCCESS] Setup complete. ROI maximized." -ForegroundColor Green
    }
    
    "build" {
        Write-P8-Header "INITIATING BUILD: LINUX PRODUCTION (DEFAULT)"
        
        Write-Host ">> Building Default Frontend (web/default)..."
        Push-Location web/default
        bun run --bun build
        Pop-Location
        
        Write-Host ">> Building Classic Frontend (web/classic)..."
        Push-Location web/classic
        bun run --bun build
        Pop-Location
        
        Write-Host ">> Compiling Linux Binary (Static, Version: $VERSION)..."
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        $env:CGO_ENABLED = "0"
        go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$VERSION'" -o new-api
        
        Write-Host "`n[SUCCESS] Build complete: 'new-api' (Linux) is ready." -ForegroundColor Green
        Write-Host ">> Target: Linux amd64 (Static Linking)" -ForegroundColor Yellow
    }

    "build-windows" {
        Write-P8-Header "INITIATING BUILD: WINDOWS CROSS-COMPILATION"
        
        Write-Host ">> Building Default Frontend (web/default)..."
        Push-Location web/default
        bun run --bun build
        Pop-Location
        
        Write-Host ">> Building Classic Frontend (web/classic)..."
        Push-Location web/classic
        bun run --bun build
        Pop-Location
        
        Write-Host ">> Compiling Windows Binary (Version: $VERSION)..."
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$VERSION'" -o new-api.exe
        
        Write-Host "`n[SUCCESS] Build complete: new-api.exe is ready." -ForegroundColor Green
    }
    
    "run" {
        Write-P8-Header "INITIATING RUN: END-TO-END CLOSURE"
        Write-Host ">> Starting Backend Dev Server..."
        go run main.go
    }
    
    "all" {
        Write-Host ">> Running all steps..." -ForegroundColor Cyan
        & $PSCommandPath setup
        & $PSCommandPath build
        & $PSCommandPath run
    }
}
