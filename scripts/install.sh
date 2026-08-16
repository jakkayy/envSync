#!/bin/sh
set -e

REPO="jakkayy/envSync"
INSTALL_DIR="/usr/local/bin"

echo "🚀 Installing envSync CLI..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$TAG" ]; then
  TAG="v0.1.0"
fi

TAR_FILE="envSync_${TAG#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$TAG/$TAR_FILE"

echo "Downloading envSync $TAG for $OS/$ARCH..."
curl -sSL "$URL" -o "/tmp/$TAR_FILE"

tar -xzf "/tmp/$TAR_FILE" -C /tmp envsync
sudo mv /tmp/envsync "$INSTALL_DIR/envsync"
rm -f "/tmp/$TAR_FILE"

echo "✅ envSync successfully installed to $INSTALL_DIR/envsync!"
envsync version
