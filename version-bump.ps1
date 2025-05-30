# NPM-style version bumping script
param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("major", "minor", "patch", "premajor", "preminor", "prepatch", "prerelease")]
    [string]$Type,
    [string]$PreID = "beta",
    [switch]$DryRun = $false
)

Write-Host "📦 NPM Version Bump Script" -ForegroundColor Green

# Function to parse semver
function Parse-SemVer {
    param([string]$Version)
    
    if ($Version -match '^v?(\d+)\.(\d+)\.(\d+)(?:-(.+))?$') {
        return @{
            Major = [int]$matches[1]
            Minor = [int]$matches[2]
            Patch = [int]$matches[3]
            PreRelease = $matches[4]
        }
    }
    return $null
}

# Function to bump version
function Bump-Version {
    param(
        [hashtable]$CurrentVersion,
        [string]$BumpType,
        [string]$PreID
    )
    
    $major = $CurrentVersion.Major
    $minor = $CurrentVersion.Minor
    $patch = $CurrentVersion.Patch
    $prerelease = $CurrentVersion.PreRelease
    
    switch ($BumpType) {
        "major" {
            $major++
            $minor = 0
            $patch = 0
            $prerelease = $null
        }
        "minor" {
            $minor++
            $patch = 0
            $prerelease = $null
        }
        "patch" {
            $patch++
            $prerelease = $null
        }
        "premajor" {
            $major++
            $minor = 0
            $patch = 0
            $prerelease = "$PreID.0"
        }
        "preminor" {
            $minor++
            $patch = 0
            $prerelease = "$PreID.0"
        }
        "prepatch" {
            $patch++
            $prerelease = "$PreID.0"
        }
        "prerelease" {
            if ($prerelease) {
                if ($prerelease -match "^$PreID\.(\d+)$") {
                    $num = [int]$matches[1] + 1
                    $prerelease = "$PreID.$num"
                } else {
                    $prerelease = "$PreID.0"
                }
            } else {
                $patch++
                $prerelease = "$PreID.0"
            }
        }
    }
    
    $newVersion = "$major.$minor.$patch"
    if ($prerelease) {
        $newVersion += "-$prerelease"
    }
    
    return $newVersion
}

# Get current version from package.json
$packageJsonPath = "npm-package/package.json"
if (-not (Test-Path $packageJsonPath)) {
    Write-Host "❌ package.json not found at $packageJsonPath" -ForegroundColor Red
    exit 1
}

$packageJson = Get-Content $packageJsonPath | ConvertFrom-Json
$currentVersion = $packageJson.version

Write-Host "Current version: $currentVersion" -ForegroundColor Cyan

# Parse current version
$parsedVersion = Parse-SemVer $currentVersion
if (-not $parsedVersion) {
    Write-Host "❌ Invalid version format in package.json: $currentVersion" -ForegroundColor Red
    exit 1
}

# Calculate new version
$newVersion = Bump-Version $parsedVersion $Type $PreID
Write-Host "New version: $newVersion" -ForegroundColor Green

if ($DryRun) {
    Write-Host "`n🔍 DRY RUN MODE - No changes will be made" -ForegroundColor Yellow
    Write-Host "Would update package.json from $currentVersion to $newVersion" -ForegroundColor Gray
    exit 0
}

# Update package.json
$packageJson.version = $newVersion
$packageJson | ConvertTo-Json -Depth 10 | Set-Content $packageJsonPath

Write-Host "`n✅ Updated package.json to version $newVersion" -ForegroundColor Green

# Ask if user wants to create release
$createRelease = Read-Host "`nCreate release tag and push? (y/N)"
if ($createRelease -eq "y" -or $createRelease -eq "Y") {
    Write-Host "`n🚀 Creating release..." -ForegroundColor Yellow
    & .\release.ps1 -Version $newVersion
} else {
    Write-Host "`n📝 Version updated. To create release later, run:" -ForegroundColor Cyan
    Write-Host "  .\release.ps1 -Version $newVersion" -ForegroundColor White
} 