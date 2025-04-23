#!/usr/bin/env bash

set -e

REPO="voltigdev/voltig"
BINARY="voltig"

# Detect latest release
LATEST_TAG=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p')
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
ASSET="${BINARY}_${LATEST_TAG#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ASSET}"

# Download to temp location
TMP_DIR=$(mktemp -d)
ARCHIVE_PATH="$TMP_DIR/$ASSET"
echo "Downloading ${ASSET} (${LATEST_TAG}) for ${OS}/${ARCH}..."
if ! curl -fsSL "$URL" -o "$ARCHIVE_PATH"; then
  echo "Failed to download archive from $URL"
  rm -rf "$TMP_DIR"
  exit 1
fi

# Extract the archive
cd "$TMP_DIR"
tar -xzf "$ARCHIVE_PATH"

# Find the binary (assume it's named voltig)
if [ ! -f "voltig" ]; then
  echo "voltig binary not found in archive."
  rm -rf "$TMP_DIR"
  exit 1
fi
chmod +x voltig

# Install to /usr/local/bin or ~/.local/bin
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  echo "No write access to /usr/local/bin, installing to $INSTALL_DIR"
fi

mv voltig "$INSTALL_DIR/$BINARY"
cd - >/dev/null
rm -rf "$TMP_DIR"

INSTALLED_PATH="$INSTALL_DIR/$BINARY"
echo "\n✅ Installed $BINARY to: $INSTALLED_PATH"

# Show which voltig is being used
WHICH_PATH=$(command -v "$BINARY" 2>/dev/null || true)
if [ -n "$WHICH_PATH" ]; then
  echo "🔎 'voltig' found in your PATH at: $WHICH_PATH"
else
  echo "⚠️  'voltig' is not currently in your PATH."
  echo "   To add it, run:"
  echo "     export PATH=\"$INSTALL_DIR:\$PATH\""
  echo "   Then restart your shell or run: source ~/.zshrc (or ~/.bashrc)"
fi

# Test install and show help
if command -v "$BINARY" >/dev/null 2>&1; then
  "$BINARY" --help || true
else
  echo "Warning: $BINARY not found in PATH. Add $INSTALL_DIR to your PATH."
fi