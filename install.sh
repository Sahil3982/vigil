#!/bin/sh
set -e

: ${VIGIL_VERSION:="latest"}
: ${VIGIL_PREFIX:="/usr/local/bin"}

ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    armv7l) ARCH=armv7 ;;
    *) echo "❌ Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
  *darwin*)   OS="darwin" ;;
  *linux*)    OS="linux" ;;
  *mingw*|*msys*) OS="windows" ;;
  *) echo "❌ Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

if [ "$VIGIL_VERSION" = "latest" ]; then
  VIGIL_VERSION=$(curl -s https://api.github.com/repos/sahil3982/vigil/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

# Fix: No space after /download/
URL="https://github.com/sahil3982/vigil/releases/download/${VIGIL_VERSION}/vigil_${VIGIL_VERSION#v}_${OS}_${ARCH}.tar.gz"

echo "📥 Downloading vigil ${VIGIL_VERSION} (${OS}/${ARCH})..."
curl -sfL "$URL" | tar -xz -C /tmp

if [ "$OS" = "windows" ]; then
  ext=".exe"
  TARGET="$HOME/bin/vigil$ext"
  mkdir -p "$HOME/bin"
  install -m 755 "/tmp/vigil$ext" "$TARGET"
  echo "✅ Installed to $TARGET"
  echo "💡 Add $HOME/bin to your PATH, or copy to a folder in PATH (e.g., C:\\Windows)."
else
  ext=""
  echo "🚚 Installing to $VIGIL_PREFIX..."
  sudo install -m 755 "/tmp/vigil$ext" "$VIGIL_PREFIX/vigil"
fi

rm -f "/tmp/vigil$ext"
echo "✅ Done! Try: vigil$ext"