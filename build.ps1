# Build script for AI Coder Context Server
# Cross-platform builds for Windows, macOS, and Linux
# Usage: .\build.ps1

param(
    [string]$Version = "1.0.0",
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"

Write-Host "Building AI Coder Context Server v$Version" -ForegroundColor Cyan

# Create output directory
if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
}

$platforms = @(
    @{GOOS = "windows"; GOARCH = "amd64"; Ext = "exe"},
    @{GOOS = "darwin"; GOARCH = "amd64"; Ext = ""},
    @{GOOS = "darwin"; GOARCH = "arm64"; Ext = ""},
    @{GOOS = "linux"; GOARCH = "amd64"; Ext = ""},
    @{GOOS = "linux"; GOARCH = "arm64"; Ext = ""}
)

foreach ($plat in $platforms) {
    $output = "$OutputDir/ai-coder-context-server-$Version-$($plat.GOOS)-$($plat.GOARCH)"
    if ($plat.Ext) {
        $output += ".$($plat.Ext)"
    }
    
    Write-Host "Building for $($plat.GOOS)/$($plat.GOARCH)..." -ForegroundColor Yellow
    
    $env:GOOS = $plat.GOOS
    $env:GOARCH = $plat.GOARCH
    
    & go build -ldflags="-s -w" -o $output ./cmd/server
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Build failed for $($plat.GOOS)/$($plat.GOARCH)" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "Build complete! Output in $OutputDir:" -ForegroundColor Green
Get-ChildItem $OutputDir | ForEach-Object { Write-Host "  $($_.Name)" }