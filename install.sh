#!/bin/sh
set -e

# Recall One-Line Installer
# Installs the latest release of Recall into ~/.local/bin/

GITHUB_REPO="vijayvenkatj/recall"
INSTALL_DIR="$HOME/.local/bin"

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    arm64|aarch64)  ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    darwin|linux) ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Get latest release from GitHub API
echo "Fetching latest release version..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to fetch latest release version. Falling back to v1.0.0"
    LATEST_RELEASE="v1.0.0"
fi

DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_RELEASE/recall_${OS}_${ARCH}.tar.gz"

echo "Downloading Recall $LATEST_RELEASE for ${OS}/${ARCH}..."
TEMP_DIR=$(mktemp -d)
curl -sL "$DOWNLOAD_URL" -o "$TEMP_DIR/recall.tar.gz"

echo "Extracting binary..."
tar -xzf "$TEMP_DIR/recall.tar.gz" -C "$TEMP_DIR"

mkdir -p "$INSTALL_DIR"
mv "$TEMP_DIR/recall" "$INSTALL_DIR/recall"
chmod +x "$INSTALL_DIR/recall"

rm -rf "$TEMP_DIR"

echo ""
echo "=============================================="
echo " Recall installed successfully to $INSTALL_DIR/recall"
echo "=============================================="
echo ""
echo "Next steps to complete the installation:"
echo "1. Ensure $INSTALL_DIR is in your PATH. If not, add this to your shell profile:"
echo "   export PATH=\"\$PATH:\$HOME/.local/bin\""
echo ""
echo "2. Run the Recall setup command:"
echo "   recall install"
