# Build script untuk npm package
# Jalankan dari root directory salat

Write-Host "Building Salat CLI for NPM distribution..." -ForegroundColor Green

# Buat struktur direktori
Write-Host "Creating directory structure..." -ForegroundColor Yellow
$binDirs = @(
    "npm-package/bin/darwin-x64",
    "npm-package/bin/darwin-arm64", 
    "npm-package/bin/linux-x64",
    "npm-package/bin/linux-arm64",
    "npm-package/bin/win32-x64",
    "npm-package/bin/win32-arm64"
)

foreach ($dir in $binDirs) {
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
}

# Build untuk semua platform
Write-Host "Building binaries for all platforms..." -ForegroundColor Yellow

# macOS
Write-Host "Building for macOS x64..." -ForegroundColor Cyan
$env:GOOS = "darwin"; $env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "npm-package/bin/darwin-x64/salat" .

Write-Host "Building for macOS ARM64..." -ForegroundColor Cyan
$env:GOOS = "darwin"; $env:GOARCH = "arm64"
go build -ldflags="-s -w" -o "npm-package/bin/darwin-arm64/salat" .

# Linux
Write-Host "Building for Linux x64..." -ForegroundColor Cyan
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "npm-package/bin/linux-x64/salat" .

Write-Host "Building for Linux ARM64..." -ForegroundColor Cyan
$env:GOOS = "linux"; $env:GOARCH = "arm64"
go build -ldflags="-s -w" -o "npm-package/bin/linux-arm64/salat" .

# Windows
Write-Host "Building for Windows x64..." -ForegroundColor Cyan
$env:GOOS = "windows"; $env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "npm-package/bin/win32-x64/salat.exe" .

Write-Host "Building for Windows ARM64..." -ForegroundColor Cyan
$env:GOOS = "windows"; $env:GOARCH = "arm64"
go build -ldflags="-s -w" -o "npm-package/bin/win32-arm64/salat.exe" .

# Reset environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

# Copy files
Write-Host "Copying package files..." -ForegroundColor Yellow
Copy-Item "README.md" "npm-package/" -Force
"MIT" | Out-File "npm-package/LICENSE" -Encoding UTF8

# Tampilkan ukuran file
Write-Host "Binary sizes:" -ForegroundColor Green
Get-ChildItem "npm-package/bin" -Recurse -File | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  $($_.FullName.Replace((Get-Location).Path + '\npm-package\bin\', '')): ${size} MB" -ForegroundColor White
}

Write-Host "Build completed! NPM package ready in npm-package/ directory" -ForegroundColor Green
Write-Host 'To test locally: cd npm-package; npm pack'    -ForegroundColor Cyan
Write-Host 'To publish:      cd npm-package; npm publish' -ForegroundColor Cyan
