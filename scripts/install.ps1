# Nordlys Installer - Windows Bootstrap
# Usage: irm https://raw.githubusercontent.com/nordlys-labs/nordlys-installer/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"
$Version = if ($env:NORDLYS_INSTALLER_VERSION) { $env:NORDLYS_INSTALLER_VERSION } else { "latest" }
$Repo = "nordlys-labs/nordlys-installer"
$InstallDir = if ($env:NORDLYS_INSTALL_DIR) { $env:NORDLYS_INSTALL_DIR } else { "$env:LOCALAPPDATA\Programs\nordlys-installer" }

function log_info { Write-Host "  $args" }
function log_success { Write-Host "  $args" }
function log_error { Write-Host "Error: $args" -ForegroundColor Red; exit 1 }

function Get-Arch {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64" -or $arch -eq "x86_64") { return "amd64" }
    if ($arch -eq "ARM64" -or $arch -eq "aarch64") { return "arm64" }
    log_error "Unsupported architecture: $arch"
}

function Get-DownloadUrl {
    $arch = Get-Arch
    if ($Version -eq "latest") {
        $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        $asset = $releases.assets | Where-Object { $_.name -eq "nordlys-installer-windows-$arch.zip" } | Select-Object -First 1
        if (-not $asset) { log_error "No release asset found for windows/$arch" }
        return $asset.browser_download_url
    }
    return "https://github.com/$Repo/releases/download/$Version/nordlys-installer-windows-$arch.zip"
}

function Install-Binary {
    $arch = Get-Arch
    $downloadUrl = Get-DownloadUrl

    log_info "Downloading nordlys-installer for windows/$arch..."
    $zipPath = "$env:TEMP\nordlys-installer.zip"

    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
    } catch {
        log_error "Failed to download: $downloadUrl"
    }

    log_success "Downloaded nordlys-installer"

    log_info "Installing to $InstallDir..."
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    Expand-Archive -Path $zipPath -DestinationPath $InstallDir -Force
    Remove-Item $zipPath -Force

    $exePath = "$InstallDir\nordlys-installer.exe"
    if (-not (Test-Path $exePath)) {
        log_error "nordlys-installer.exe not found after extraction"
    }

    log_success "Installed to $exePath"

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        log_info "Adding to user PATH..."
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
        log_success "Added $InstallDir to PATH. Restart your terminal for changes to take effect."
    }
}

Write-Host ""
Write-Host "=========================================="
Write-Host "  Nordlys Installer Bootstrap"
Write-Host "=========================================="
Write-Host ""

Install-Binary

Write-Host ""
log_success "Installation complete!"
Write-Host ""
Write-Host "Get started:"
Write-Host "   nordlys-installer              # Interactive mode"
Write-Host "   nordlys-installer list         # List supported tools"
Write-Host "   nordlys-installer --help       # View all options"
Write-Host ""
