#!/bin/bash

set -e

REPO="RamiroCyber/go-api-gen"
VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")
BINARY_NAME="go-api-gen"
INSTALL_DIR="$HOME/.local/bin"
BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Validate GitHub API response
if [ -z "$VERSION" ]; then
  echo "❌ Failed to fetch latest release version from GitHub. Check your internet connection or repository settings."
  exit 1
fi

# Determine OS and architecture
case "$OS" in
  linux)
    case "$ARCH" in
      x86_64) BINARY_FILE="${BINARY_NAME}-linux-amd64" ;;
      aarch64) BINARY_FILE="${BINARY_NAME}-linux-arm64" ;;
      *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    ;;
  darwin)
    case "$ARCH" in
      x86_64) BINARY_FILE="${BINARY_NAME}-macos-amd64" ;;
      arm64) BINARY_FILE="${BINARY_NAME}-macos-arm64" ;;
      *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    ;;
  *)
    echo "❌ Unsupported operating system: $OS"
    exit 1
    ;;
esac

# Create installation directory
mkdir -p "$INSTALL_DIR"

# Download the binary
echo "⬇️ Downloading $BINARY_FILE (version $VERSION)..."
curl -L -o "$BINARY_PATH" "https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_FILE}" || {
  echo "❌ Failed to download binary. Ensure the release exists at https://github.com/${REPO}/releases/tag/${VERSION}"
  exit 1
}

# Make binary executable
chmod +x "$BINARY_PATH"

# Check if INSTALL_DIR is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "⚠️ $INSTALL_DIR is not in your PATH. Adding it to ~/.bashrc or ~/.zshrc..."
  SHELL_CONFIG=""
  if [ -n "$ZSH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
  elif [ -n "$BASH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
  else
    echo "⚠️ Could not detect shell. Please manually add to your shell config:"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
  fi
  if [ -n "$SHELL_CONFIG" ]; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_CONFIG"
    echo "✅ Added $INSTALL_DIR to $SHELL_CONFIG. Run 'source $SHELL_CONFIG' to apply changes."
  fi
else
  echo "✅ $INSTALL_DIR is already in PATH"
fi

# Test the binary
echo "ℹ️ Testing binary version:"
if ! "$BINARY_PATH" --version; then
  echo "⚠️ The binary does not implement --version yet. Try running:"
  echo "  $BINARY_PATH generate module test"
  echo "If templates fail to load, ensure the binary was built with embedded templates."
fi

echo "✅ Installation completed: $BINARY_PATH"