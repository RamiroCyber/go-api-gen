#!/bin/bash

set -e

REPO="RamiroCyber/go-api-gen"
VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
BINARY_NAME="go-api-gen"
INSTALL_DIR="$HOME/.local/bin"
BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"

OS=$(uname | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)
    BINARY_FILE="${BINARY_NAME}-linux"
    ;;
  darwin)
    BINARY_FILE="${BINARY_NAME}-macos"
    ;;
  *)
    echo "❌ Sistema operacional não suportado: $OS"
    exit 1
    ;;
esac

mkdir -p "$INSTALL_DIR"

echo "⬇️  Baixando $BINARY_FILE..."
curl -L -o "$BINARY_PATH" "https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_FILE}"

chmod +x "$BINARY_PATH"

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "⚠️  Adicione '$INSTALL_DIR' ao seu PATH:"
  echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
else
  echo "✅ '$INSTALL_DIR' já está no PATH"
fi

echo "✅ Instalação concluída: $BINARY_PATH"
echo "ℹ️  Testando versão:"
"$BINARY_PATH" --version || echo "⚠️  O binário ainda não implementa --version"
