#!/bin/sh
set -e

: ${VIGIL_VERSION:="latest"}
: ${VIGIL_PREFIX:="/usr/local/bin"}

ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    armv7l) ARCH=armv7 ;;
    *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$OS" = "darwin" ]; then
    OS="darwin"
elif [ "$OS" = "linux" ]; then
    OS="linux"
else
    echo "Unsupported OS: $OS" >&2; exit 1
fi

if [ "$VIGIL_VERSION" = "latest" ]; then
    VIGIL_VERSION=$(curl -s https://api.github.com/repos/sahil3982/vigil/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

URL="https://github.com/sahil3982/vigil/releases/download/${VIGIL_VERSION}/vigil_${VIGIL_VERSION#v}_${OS}_${ARCH}.tar.gz"

echo "📥 Downloading vigil ${VIGIL_VERSION} (${OS}/${ARCH})..."
curl -sfL "$URL" | tar -xz -C /tmp

echo "🚚 Installing to $VIGIL_PREFIX..."
sudo install -m 755 /tmp/vigil "$VIGIL_PREFIX/vigil"

rm -f /tmp/vigil

echo "✅ Done! Try: vigil cpu"