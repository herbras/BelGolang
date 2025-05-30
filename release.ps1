# Script untuk menjalankan GoReleaser secara lokal
# Pastikan GoReleaser sudah terinstall: https://goreleaser.com/install/

$goreleaserInstalled = $null
try {
    $goreleaserInstalled = Get-Command goreleaser -ErrorAction SilentlyContinue
} catch {}

if ($null -eq $goreleaserInstalled) {
    Write-Host "❌ GoReleaser tidak ditemukan. Silakan install terlebih dahulu:" -ForegroundColor Red
    Write-Host "   https://goreleaser.com/install/" -ForegroundColor Yellow
    Write-Host "   Atau gunakan: go install github.com/goreleaser/goreleaser@latest" -ForegroundColor Yellow
    exit 1
}

$status = git status --porcelain
if ($status) {
    Write-Host "⚠️ Ada perubahan yang belum di-commit:" -ForegroundColor Yellow
    git status --short
    
    $confirmation = Read-Host "Lanjutkan release? (y/N)"
    if ($confirmation -ne "y") {
        Write-Host "Release dibatalkan." -ForegroundColor Red
        exit 1
    }
}

$currentVersion = git describe --tags --abbrev=0 2>$null
if ($LASTEXITCODE -ne 0) {
    $currentVersion = "v0.0.0"
}

Write-Host "📦 Versi saat ini: $currentVersion" -ForegroundColor Cyan

$newVersion = Read-Host "Masukkan versi baru (contoh: v1.0.0)"
if (-not $newVersion) {
    Write-Host "Release dibatalkan." -ForegroundColor Red
    exit 1
}

Write-Host "\n🚀 Akan merilis versi $newVersion" -ForegroundColor Green
$confirmation = Read-Host "Lanjutkan? (y/N)"
if ($confirmation -ne "y") {
    Write-Host "Release dibatalkan." -ForegroundColor Red
    exit 1
}

Write-Host "\n📝 Membuat tag $newVersion..." -ForegroundColor Yellow
git tag -a $newVersion -m "Release $newVersion"

Write-Host "\n🚀 Menjalankan GoReleaser..." -ForegroundColor Yellow
goreleaser release --clean --skip-publish

Write-Host "\n✅ Build selesai! File release tersedia di folder ./dist" -ForegroundColor Green
Write-Host "\n📋 Untuk publish ke npm dan GitHub:" -ForegroundColor Cyan
Write-Host "   1. Push tag: git push origin $newVersion" -ForegroundColor White
Write-Host "   2. Jalankan GitHub Actions workflow" -ForegroundColor White