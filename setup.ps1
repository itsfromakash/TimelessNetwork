# PowerShell Setup & Build Automation Script
Write-Host "=== Setting up Native Build Environment ===" -ForegroundColor Cyan

# Build ෆෝල්ඩරයක් නොමැති නම් සාදා ගැනීම
if (!(Test-Path "build")) {
    New-Item -ItemType Directory -Path "build" | Out-Null
}

Set-Location build

Write-Host "Running CMake Configuration..." -ForegroundColor Yellow
cmake -G "Visual Studio 17 2022" -A x64 ..

Write-Host "Building Project Files..." -ForegroundColor Yellow
cmake --build . --config Release

Write-Host "=== Build Completed Successfully! Executing MainApp ===" -ForegroundColor Green

# Build වූ Executable එක වෙත ගොස් එය ක්‍රියාත්මක කිරීම
Set-Location Release
if (Test-Path "MainApp.exe") {
    .\MainApp.exe
} else {
    Write-Host "Error: Executable not found." -ForegroundColor Red
}

Set-Location ../..