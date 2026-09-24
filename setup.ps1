# 1. Folders සාදා ගැනීම
New-Item -ItemType Directory -Force -Path "include", "lib", "build" | Out-Null

# 2. Go Code එක C Static Library එකක් ලෙස Compile කිරීම
Write-Host "[BUILD] Compiling Go P2P Engine to .a static library..." -ForegroundColor Cyan
Set-Location src/P2P
go build -buildmode=c-archive -o "../../lib/p2p_engine.a" p2p.go
Move-Item -Path "../../lib/p2p_engine.h" -Destination "../../include/p2p_engine.h" -Force
Set-Location ../..

# 3. CMake හරහා C++ Executable එක Compile කිරීම
Write-Host "[BUILD] Building C++ Core Engine executable..." -ForegroundColor Cyan
Set-Location build
cmake -G "MinGW Makefiles" ..
cmake --build .
Set-Location ..

Write-Host "`n[SUCCESS] Build Completed Successfully!" -ForegroundColor Green
Write-Host "Run command: .\build\TimelessNetworkApp.exe" -ForegroundColor Yellow