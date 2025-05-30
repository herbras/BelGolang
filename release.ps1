# Release script - Creates new version tag and triggers CI/CD
param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    [switch]$DryRun = $false
)

Write-Host "🚀 Salat CLI Release Script" -ForegroundColor Green

# Validate version format
if ($Version -notmatch '^v?\d+\.\d+\.\d+$') {
    Write-Host "❌ Invalid version format. Use format: 1.6.2 or v1.6.2" -ForegroundColor Red
    exit 1
}

# Normalize version (remove v prefix if present)
$CleanVersion = $Version -replace '^v', ''
$TagVersion = "v$CleanVersion"

Write-Host "Version: $CleanVersion" -ForegroundColor Cyan
Write-Host "Git Tag: $TagVersion" -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "`n🔍 DRY RUN MODE - No changes will be made" -ForegroundColor Yellow
}

# Check if working directory is clean
$gitStatus = & git status --porcelain
if ($gitStatus -and -not $DryRun) {
    Write-Host "❌ Working directory is not clean. Please commit or stash changes." -ForegroundColor Red
    Write-Host "Uncommitted changes:" -ForegroundColor Yellow
    $gitStatus | ForEach-Object { Write-Host "  $_" -ForegroundColor Red }
    exit 1
}

# Check if tag already exists
$existingTag = & git tag -l $TagVersion
if ($existingTag) {
    Write-Host "❌ Tag $TagVersion already exists!" -ForegroundColor Red
    exit 1
}

# Update package.json version locally (for reference)
if (-not $DryRun) {
    Write-Host "`n📝 Updating npm-package/package.json..." -ForegroundColor Yellow
    $packageJsonPath = "npm-package/package.json"
    
    if (Test-Path $packageJsonPath) {
        $packageJson = Get-Content $packageJsonPath | ConvertFrom-Json
        $packageJson.version = $CleanVersion
        $packageJson | ConvertTo-Json -Depth 10 | Set-Content $packageJsonPath
        Write-Host "✅ Updated package.json to version $CleanVersion" -ForegroundColor Green
    } else {
        Write-Host "⚠️  package.json not found at $packageJsonPath" -ForegroundColor Yellow
    }
}

# Create and push tag
if (-not $DryRun) {
    Write-Host "`n🏷️  Creating git tag..." -ForegroundColor Yellow
    
    # Create annotated tag with release notes
    $tagMessage = @"
Release $TagVersion

🎯 Features:
- Ultra-optimized binaries (79% size reduction)
- Embedded timezone database
- Cross-platform support (Windows, macOS, Linux)
- ARM64 and AMD64 architectures
- WASM support for browsers

📦 Installation:
npm install -g salat-cli@$CleanVersion

🔗 GitHub: https://github.com/herbras/BelGolang
📚 Docs: https://github.com/herbras/BelGolang#readme
"@
    
    & git add .
    & git commit -m "chore: release $TagVersion" --allow-empty
    & git tag -a $TagVersion -m $tagMessage
    
    Write-Host "✅ Created tag $TagVersion" -ForegroundColor Green
    
    # Push changes and tags
    Write-Host "`n📤 Pushing to remote..." -ForegroundColor Yellow
    & git push origin main
    & git push origin $TagVersion
    
    Write-Host "✅ Pushed tag $TagVersion to remote" -ForegroundColor Green
    
    Write-Host "`n🎉 Release process started!" -ForegroundColor Green
    Write-Host "📋 What happens next:" -ForegroundColor Cyan
    Write-Host "  1. GitHub Actions will build binaries for all platforms" -ForegroundColor White
    Write-Host "  2. Binaries will be optimized (CGO_ENABLED=0 + UPX-ready)" -ForegroundColor White
    Write-Host "  3. NPM package will be published automatically" -ForegroundColor White
    Write-Host "  4. GitHub release will be created with artifacts" -ForegroundColor White
    
    Write-Host "`n🔗 Monitor progress:" -ForegroundColor Cyan
    Write-Host "  GitHub Actions: https://github.com/herbras/BelGolang/actions" -ForegroundColor Blue
    Write-Host "  NPM Package: https://www.npmjs.com/package/salat-cli" -ForegroundColor Blue
    Write-Host "  Releases: https://github.com/herbras/BelGolang/releases" -ForegroundColor Blue
    
} else {
    Write-Host "`n🔍 DRY RUN - Would have created tag: $TagVersion" -ForegroundColor Yellow
    Write-Host "Command to run: git tag -a $TagVersion -m 'Release $TagVersion'" -ForegroundColor Gray
}

Write-Host "`n✨ Release script completed!" -ForegroundColor Green