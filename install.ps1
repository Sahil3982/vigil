# vigil install script for Windows PowerShell
param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:USERPROFILE\.vigil"
)

$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $Arch = "arm64"
}

if ($Version -eq "latest") {
    Write-Host "🔍 Fetching latest version..." -ForegroundColor Cyan
    $ReleaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/sahil3982/vigil/releases/latest"
    $Version = $ReleaseInfo.tag_name
    if (-not $Version) {
        Write-Error "❌ No release found"
        exit 1
    }
}

$Version = $Version.TrimStart('v')
$Url = "https://github.com/sahil3982/vigil/releases/download/v$Version/vigil_v$Version_windows_$Arch.tar.gz"

Write-Host "📥 Downloading vigil v$Version for Windows/$Arch..." -ForegroundColor Cyan

$TempDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }
$ZipPath = Join-Path $TempDir "vigil.tar.gz"

try {
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath

    # ✅ FIX: Use PowerShell native extractor
    $ExtractDir = Join-Path $TempDir "extract"
    New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null

    Expand-Archive -Path $ZipPath -DestinationPath $ExtractDir -Force

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

    $ExePath = Get-ChildItem -Path $ExtractDir -Filter "vigil.exe" | Select-Object -First 1
    if ($ExePath) {
        $TargetPath = Join-Path $InstallDir "vigil.exe"
        Move-Item -Path $ExePath.FullName -Destination $TargetPath -Force
        Write-Host "✅ Installed to: $TargetPath" -ForegroundColor Green
    } else {
        Write-Error "❌ vigil.exe not found in archive"
        exit 1
    }

    $CurrentPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    if ($CurrentPath -notlike "*$InstallDir*") {
        [System.Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
        Write-Host "🔧 Added to PATH (restart terminal to apply)" -ForegroundColor Yellow
    } else {
        Write-Host "💡 vigil is already in PATH" -ForegroundColor Green
    }

    Write-Host "🎉 Done! Run 'vigil --help' to get started." -ForegroundColor Green
}
catch {
    Write-Error "❌ Installation failed: $_"
    exit 1
}
finally {
    Remove-Item $TempDir -Recurse -Force -ErrorAction SilentlyContinue
}