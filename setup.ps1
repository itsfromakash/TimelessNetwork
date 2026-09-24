# Save the current working directory as the root folder
$RootDir = Get-Location

# 1. Clean Build Directories
Write-Host "[CLEAN] Cleaning old build files..." -ForegroundColor Yellow

# Delete the old build directory if it exists
if (Test-Path "$RootDir\build") { Remove-Item -Recurse -Force "$RootDir\build" }

# Create fresh 'include', 'lib', and 'build' directories silently
New-Item -ItemType Directory -Force -Path "include", "lib", "build" | Out-Null

# 2. Navigate to the exact location of the Go P2P file
$GoDir = "$RootDir\src\TimelessNetwork"
Write-Host "[BUILD] Navigating to Go directory: $GoDir" -ForegroundColor Cyan

Set-Location $GoDir

# 3. Build the Go Static Library
Write-Host "[BUILD] Compiling Go P2P Engine to static library..." -ForegroundColor Cyan

# Compile Go code into a C-compatible static archive library (.a file)
go build -buildmode=c-archive -o "$RootDir\lib\p2p_engine.a" main.go

# Move the auto-generated C header file (.h) into the 'include' directory
if (Test-Path "$RootDir\lib\p2p_engine.h") {
    Move-Item -Path "$RootDir\lib\p2p_engine.h" -Destination "$RootDir\include\p2p_engine.h" -Force
}

# 4. Return to the root directory and compile the C++ Executable using CMake
Set-Location "$RootDir\build"
Write-Host "[BUILD] Building C++ Core Engine executable..." -ForegroundColor Cyan

# Generate MinGW build files from the parent CMakeLists.txt
cmake -G "MinGW Makefiles" ".."
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] CMake configuration failed!" -ForegroundColor Red
    Set-Location $RootDir
    return
}

# Compile the C++ project using CMake
cmake --build .
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Compilation failed!" -ForegroundColor Red
    Set-Location $RootDir
    return
}

# Return back to the original root folder
Set-Location $RootDir

Write-Host "`n[SUCCESS] Build Completed Successfully!" -ForegroundColor Green
Write-Host "Run command: .\build\TimelessNetworkApp.exe" -ForegroundColor Yellow
