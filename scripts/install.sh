#!/usr/bin/env bash

set -e

REPO="voltigdev/voltig"
BINARY="voltig"

# Detect latest release
LATEST_TAG=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep -oP '"tag_name": "\K(.*)(?=")')
if [ -z "$LATEST_TAG" ]; then
  echo "Could not determine latest release. Exiting."
  exit 1
fi

# Detect OS and ARCH
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  armv8*|aarch64) ARCH="arm64" ;;
  armv7*) ARCH="armv7" ;;
esac

# Compose download URL
ASSET="${BINARY}-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ASSET}"

# Download to temp location
TMP=$(mktemp)
echo "Downloading ${BINARY} (${LATEST_TAG}) for ${OS}/${ARCH}..."
if ! curl -fsSL "$URL" -o "$TMP"; then
  echo "Failed to download binary from $URL"
  exit 1
fi

chmod +x "$TMP"

# Install to /usr/local/bin or ~/.local/bin
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  echo "No write access to /usr/local/bin, installing to $INSTALL_DIR"
fi

mv "$TMP" "$INSTALL_DIR/$BINARY"

echo "Installed $BINARY to $INSTALL_DIR/$BINARY"
echo "Make sure $INSTALL_DIR is in your PATH."

# Test install
if command -v "$BINARY" >/dev/null 2>&1; then
  "$BINARY" --help || true
else
  echo "Warning: $BINARY not found in PATH. Add $INSTALL_DIR to your PATH."
fi