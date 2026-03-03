#!/bin/bash
set -euo pipefail

VERSION="${NORDLYS_INSTALLER_VERSION:-latest}"
REPO="nordlys-labs/nordlys-installer"
INSTALL_DIR="${NORDLYS_INSTALL_DIR:-/usr/local/bin}"

log_info() { echo "🔹 $*"; }
log_success() { echo "✅ $*"; }
log_error() { echo "❌ $*" >&2; }

detect_platform() {
	OS=$(uname -s | tr '[:upper:]' '[:lower:]')
	ARCH=$(uname -m)
	
	case "$ARCH" in
		x86_64) ARCH="amd64" ;;
		aarch64|arm64) ARCH="arm64" ;;
		*) log_error "Unsupported architecture: $ARCH"; exit 1 ;;
	esac
	
	case "$OS" in
		linux|darwin) ;;
		*) log_error "Unsupported OS: $OS"; exit 1 ;;
	esac
}

download_installer() {
	local download_url
	
	if [ "$VERSION" = "latest" ]; then
		download_url="https://github.com/$REPO/releases/latest/download/nordlys-installer-${OS}-${ARCH}"
	else
		download_url="https://github.com/$REPO/releases/download/${VERSION}/nordlys-installer-${OS}-${ARCH}"
	fi
	
	log_info "Downloading nordlys-installer for $OS/$ARCH..."
	
	if ! curl -fsSL "$download_url" -o /tmp/nordlys-installer; then
		log_error "Failed to download installer from $download_url"
		exit 1
	fi
	
	chmod +x /tmp/nordlys-installer
	log_success "Downloaded nordlys-installer"
}

install_binary() {
	log_info "Installing to $INSTALL_DIR..."
	
	if [ ! -w "$INSTALL_DIR" ]; then
		log_info "Requires sudo for installation to $INSTALL_DIR"
		sudo mv /tmp/nordlys-installer "$INSTALL_DIR/nordlys-installer"
	else
		mv /tmp/nordlys-installer "$INSTALL_DIR/nordlys-installer"
	fi
	
	log_success "Installed to $INSTALL_DIR/nordlys-installer"
}

main() {
	echo "=========================================="
	echo "  Nordlys Installer Bootstrap"
	echo "=========================================="
	echo ""
	
	detect_platform
	download_installer
	install_binary
	
	echo ""
	log_success "Installation complete!"
	echo ""
	echo "🚀 Get started:"
	echo "   nordlys-installer              # Interactive mode"
	echo "   nordlys-installer list         # List supported tools"
	echo "   nordlys-installer --help       # View all options"
	echo ""
}

main "$@"
