# PowerShell script to build and install a Go project

# Exit on any error
$ErrorActionPreference = "Stop"

# Define paths
$binDir = "bin"
$exeName = "key.exe"
$installDir = "$env:USERPROFILE\.key"
$installPath = "$installDir\$exeName"

Write-Host "Building Key..." -ForegroundColor Cyan

# Create bin directory if it doesn't exist
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

# Build the Go project
# -ldflags="-s -w" strips debug info and reduces binary size
# -trimpath removes local path info from the binary
go build -ldflags="-s -w" -trimpath -o "$binDir\$exeName"

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Build successful!" -ForegroundColor Green

# Create installation directory if it doesn't exist
Write-Host "Creating installation directory: $installDir" -ForegroundColor Cyan
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Move the executable to the installation directory
Write-Host "Moving $exeName to $installPath" -ForegroundColor Cyan
Copy-Item -Path "$binDir\$exeName" -Destination $installPath -Force

# Add to user PATH if not already present
Write-Host "Updating PATH environment variable..." -ForegroundColor Cyan
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")

if ($currentPath -notlike "*$installDir*") {
    $newPath = $currentPath.TrimEnd(';') + ";$installDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "Added $installDir to user PATH" -ForegroundColor Green
    Write-Host "Please restart your terminal or run 'refreshenv' to use the updated PATH" -ForegroundColor Yellow
} else {
    Write-Host "$installDir is already in PATH" -ForegroundColor Yellow
}

Write-Host "Installation complete!" -ForegroundColor Green