$RootDir = Get-Location

# 1. Clean Build Directories
Write-Host "[CLEAN] Cleaning old build files..." -ForegroundColor Yellow
if (Test-Path "$RootDir\build") { Remove-Item -Recurse -Force "$RootDir\build" }
New-Item -ItemType Directory -Force -Path "include", "lib", "build" | Out-Null

# 2. Go P2P File එක ඇති exact location එකට යාම
$GoDir = "$RootDir\src\TimelessNetwork"
Write-Host "[BUILD] Navigating to Go directory: $GoDir" -ForegroundColor Cyan

Set-Location $GoDir

# 3. Go Static Library එක Build කිරීම
Write-Host "[BUILD] Compiling Go P2P Engine to static library..." -ForegroundColor Cyan
go build -buildmode=c-archive -o "$RootDir\lib\p2p_engine.a" main.go

if (Test-Path "$RootDir\lib\p2p_engine.h") {
    Move-Item -Path "$RootDir\lib\p2p_engine.h" -Destination "$RootDir\include\p2p_engine.h" -Force
}

# 4. Root එකට පැමිණ CMake හරහා C++ Executable එක Compile කිරීම
Set-Location "$RootDir\build"
Write-Host "[BUILD] Building C++ Core Engine executable..." -ForegroundColor Cyan

cmake -G "MinGW Makefiles" ".."
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] CMake configuration failed!" -ForegroundColor Red
    Set-Location $RootDir
    return
}

cmake --build .
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Compilation failed!" -ForegroundColor Red
    Set-Location $RootDir
    return
}

Set-Location $RootDir

Write-Host "`n[SUCCESS] Build Completed Successfully!" -ForegroundColor Green
Write-Host "Run command: .\build\TimelessNetworkApp.exe" -ForegroundColor Yellow