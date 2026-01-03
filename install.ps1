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
    $ReleaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/sahil3982/vigil/releases/latest" -ErrorAction Stop
    $Version = $ReleaseInfo.tag_name
    if (-not $Version) {
        Write-Error "❌ No release found"
        exit 1
    }
}

# Remove 'v' prefix for URL construction
$VersionNumber = $Version.TrimStart('v')

Write-Host "📥 Downloading vigil $Version for Windows/$Arch..." -ForegroundColor Cyan

$TempDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }

try {
    # Try two possible URL patterns
    $Urls = @(
        "https://github.com/sahil3982/vigil/releases/download/$Version/vigil_${VersionNumber}_windows_$Arch.tar.gz",
        "https://github.com/sahil3982/vigil/releases/download/$Version/vigil_v${VersionNumber}_windows_$Arch.tar.gz"
    )
    
    $success = $false
    foreach ($Url in $Urls) {
        try {
            Write-Host "Trying: $Url" -ForegroundColor DarkGray
            $ZipPath = Join-Path $TempDir "vigil.tar.gz"
            Invoke-WebRequest -Uri $Url -OutFile $ZipPath -ErrorAction Stop
            
            $ExtractDir = Join-Path $TempDir "extract"
            New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null
            
            # Use 7zip or tar if available, fallback to .NET
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                tar -xzf $ZipPath -C $ExtractDir
            } else {
                # Manual extraction for .tar.gz
                $gzipStream = [System.IO.Compression.GZipStream]::new(
                    [System.IO.File]::OpenRead($ZipPath),
                    [System.IO.Compression.CompressionMode]::Decompress
                )
                $tarPath = Join-Path $TempDir "vigil.tar"
                $tarFile = [System.IO.File]::Create($tarPath)
                $gzipStream.CopyTo($tarFile)
                $gzipStream.Close()
                $tarFile.Close()
                
                # Extract tar
                tar -xf $tarPath -C $ExtractDir
            }
            
            $success = $true
            break
        }
        catch {
            Write-Host "Failed: $_" -ForegroundColor DarkYellow
            continue
        }
    }
    
    if (-not $success) {
        Write-Error "❌ Failed to download from any URL pattern"
        exit 1
    }

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

    $ExePath = Get-ChildItem -Path $ExtractDir -Recurse -Filter "vigil.exe" | Select-Object -First 1
    if ($ExePath) {
        $TargetPath = Join-Path $InstallDir "vigil.exe"
        Move-Item -Path $ExePath.FullName -Destination $TargetPath -Force
        Write-Host "✅ Installed to: $TargetPath" -ForegroundColor Green
    } else {
        Write-Error "❌ vigil.exe not found in archive"
        exit 1
    }

    # Add to PATH
    $CurrentPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    if ($CurrentPath -notlike "*$InstallDir*") {
        [System.Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
        Write-Host "🔧 Added to PATH (restart terminal to apply)" -ForegroundColor Yellow
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