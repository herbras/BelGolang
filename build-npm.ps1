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

# Build for all platforms with optimization flags
# Note: CGO_ENABLED=0 creates static binaries with embedded timezone data
$platforms = @(
    @{GOOS="windows"; GOARCH="amd64"; EXT=".exe"; DIR="win32-x64"}
    @{GOOS="windows"; GOARCH="arm64"; EXT=".exe"; DIR="win32-arm64"}
    @{GOOS="darwin"; GOARCH="amd64"; EXT=""; DIR="darwin-x64"}
    @{GOOS="darwin"; GOARCH="arm64"; EXT=""; DIR="darwin-arm64"}
    @{GOOS="linux"; GOARCH="amd64"; EXT=""; DIR="linux-x64"}
    @{GOOS="linux"; GOARCH="arm64"; EXT=""; DIR="linux-arm64"}
)

# Get version for ldflags
$version = & git describe --tags --always 2>$null
if (-not $version) { $version = "dev" }

foreach ($platform in $platforms) {
    Write-Host "Building for $($platform.DIR)..." -ForegroundColor Green
    
    $env:GOOS = $platform.GOOS
    $env:GOARCH = $platform.GOARCH
    $env:CGO_ENABLED = "0"  # Static binary with embedded timezone database
    
    $outputPath = "npm-package/bin/$($platform.DIR)/salat$($platform.EXT)"
    
    # Build with optimization flags - fixed PowerShell quoting
    $ldflags = "-s -w -X main.version=$version"
    & go build -ldflags $ldflags -trimpath -o $outputPath .
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Built $($platform.DIR)" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to build $($platform.DIR)" -ForegroundColor Red
        exit 1
    }
}

# Reset environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue

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

# UPX Compression option
Write-Host "`n🗜️  UPX Compression:" -ForegroundColor Cyan
Write-Host "Note: CI/CD automatically applies UPX compression to published NPM packages" -ForegroundColor Gray
$upxChoice = Read-Host "Apply UPX compression to local builds? (Y/n)"

if ($upxChoice -ne "n" -and $upxChoice -ne "N") {
    # Check if UPX is available
    $upxPath = Get-Command upx -ErrorAction SilentlyContinue
    if ($upxPath) {
        Write-Host "Applying UPX compression to all binaries..." -ForegroundColor Yellow
        
        $compressedCount = 0
        $totalOriginalSize = 0
        $totalCompressedSize = 0
        
        Get-ChildItem "npm-package/bin" -Recurse -File -Include "salat*" | ForEach-Object {
            $originalSize = $_.Length
            $totalOriginalSize += $originalSize
            
            Write-Host "Compressing $($_.Name)..." -ForegroundColor Gray
            & upx --best --lzma $_.FullName 2>$null
            
            if ($LASTEXITCODE -eq 0) {
                $compressedSize = (Get-Item $_.FullName).Length
                $totalCompressedSize += $compressedSize
                $reduction = [math]::Round(($originalSize - $compressedSize) / $originalSize * 100, 1)
                Write-Host "  ✅ $($_.Name): ${reduction}% reduction" -ForegroundColor Green
                $compressedCount++
            } else {
                Write-Host "  ⚠️ Failed to compress $($_.Name)" -ForegroundColor Red
                $totalCompressedSize += $originalSize
            }
        }
        
        if ($compressedCount -gt 0) {
            Write-Host "`n📊 UPX Compression Summary:" -ForegroundColor Cyan
            Write-Host "  Compressed files: $compressedCount" -ForegroundColor White
            $totalReduction = [math]::Round(($totalOriginalSize - $totalCompressedSize) / $totalOriginalSize * 100, 1)
            Write-Host "  Total size reduction: ${totalReduction}%" -ForegroundColor Yellow
            
            Write-Host "`nCompressed binary sizes:" -ForegroundColor Green
            Get-ChildItem "npm-package/bin" -Recurse -File | ForEach-Object {
                $size = [math]::Round($_.Length / 1MB, 2)
                Write-Host "  $($_.FullName.Replace((Get-Location).Path + '\npm-package\bin\', '')): ${size} MB" -ForegroundColor Magenta
            }
            
            # Test one binary to ensure it works
            Write-Host "`n🧪 Testing compressed binaries..." -ForegroundColor Cyan
            $testBinary = Get-ChildItem "npm-package/bin" -Recurse -File -Include "salat*" | Select-Object -First 1
            if ($testBinary) {
                $testResult = & $testBinary.FullName --help 2>$null
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "✅ Compressed binaries work correctly!" -ForegroundColor Green
                } else {
                    Write-Host "⚠️  Compressed binaries may have issues" -ForegroundColor Red
                }
            }
        }
    } else {
        Write-Host "⚠️  UPX not found. Install with: winget install upx" -ForegroundColor Red
        Write-Host "Continuing without compression..." -ForegroundColor Gray
    }
} else {
    Write-Host "Skipping UPX compression..." -ForegroundColor Gray
}

Write-Host "`nBuild completed! NPM package ready in npm-package/ directory" -ForegroundColor Green
Write-Host 'To test locally: cd npm-package; npm pack'    -ForegroundColor Cyan
Write-Host 'To publish:      cd npm-package; npm publish' -ForegroundColor Cyan
