#!/bin/bash

set -e

# Install script for tusk CLI
# Downloads binary from GitHub releases and installs it

REPO="sboy99/projects.go"
BINARY_NAME="tusk"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)
            echo "linux"
            ;;
        Darwin*)
            echo "darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Detect architecture
detect_arch() {
    local arch
    arch=$(uname -m)
    
    case "$arch" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    echo "$version"
}

# Download and install
install_tusk() {
    local os arch version download_url archive_name binary_name install_path
    
    os=$(detect_os)
    arch=$(detect_arch)
    
    echo -e "${GREEN}Detected system: ${os}/${arch}${NC}"
    
    # Check if architecture is supported
    case "${os}/${arch}" in
        linux/amd64|linux/arm64|darwin/amd64|darwin/arm64)
            ;;
        windows/amd64)
            echo -e "${YELLOW}Windows installation is not fully supported in this script.${NC}"
            echo -e "${YELLOW}Please download the .zip file manually from GitHub releases.${NC}"
            exit 1
            ;;
        *)
            echo -e "${RED}Error: Currently unsupported architecture: ${os}/${arch}${NC}"
            echo -e "${YELLOW}Supported architectures:${NC}"
            echo "  - linux/amd64"
            echo "  - linux/arm64"
            echo "  - darwin/amd64"
            echo "  - darwin/arm64"
            exit 1
            ;;
    esac
    
    # Get version
    if [ "$VERSION" = "latest" ]; then
        echo "Fetching latest version..."
        version=$(get_latest_version)
        if [ -z "$version" ]; then
            echo -e "${RED}Error: Could not fetch latest version${NC}"
            exit 1
        fi
    else
        version="$VERSION"
    fi
    
    echo -e "${GREEN}Installing tusk ${version}...${NC}"
    
    # Set file names
    archive_name="tusk-${version}-${os}-${arch}.tar.gz"
    binary_name="$BINARY_NAME"
    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"
    
    # Create temporary directory
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    # Download archive
    echo "Downloading ${archive_name}..."
    if ! curl -fsSL "$download_url" -o "${tmp_dir}/${archive_name}"; then
        echo -e "${RED}Error: Failed to download ${archive_name}${NC}"
        echo -e "${YELLOW}Available releases: https://github.com/${REPO}/releases${NC}"
        exit 1
    fi
    
    # Extract archive
    echo "Extracting archive..."
    cd "$tmp_dir"
    tar -xzf "$archive_name"
    
    # Find the binary in the extracted directory
    binary_path=$(find . -name "$binary_name" -type f | head -n 1)
    if [ -z "$binary_path" ]; then
        echo -e "${RED}Error: Binary not found in archive${NC}"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$binary_path"
    
    # Install binary
    install_path="${INSTALL_DIR}/${binary_name}"
    echo "Installing to ${install_path}..."
    
    # Check if we need sudo
    if [ ! -w "$INSTALL_DIR" ]; then
        echo "Requires sudo to install to ${INSTALL_DIR}"
        sudo mv "$binary_path" "$install_path"
        sudo chown root:root "$install_path" 2>/dev/null || true
    else
        mv "$binary_path" "$install_path"
    fi
    
    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Successfully installed tusk ${version}${NC}"
        echo ""
        echo "Run 'tusk version' to verify installation"
    else
        echo -e "${YELLOW}Warning: tusk installed but not found in PATH${NC}"
        echo "You may need to add ${INSTALL_DIR} to your PATH"
    fi
}

# Main
main() {
    echo "Tusk CLI Installer"
    echo "=================="
    echo ""
    
    install_tusk
}

# Run main function
main "$@"

