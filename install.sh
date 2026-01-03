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

UNAME_S=$(uname -s)
case $UNAME_S in
    *Darwin*)   OS="darwin" ;;
    *Linux*)    OS="linux" ;;
    *MINGW*|*MSYS*|*CYGWIN*) OS="windows" ;;
    *) echo "❌ Unsupported OS: $UNAME_S" >&2; exit 1 ;;
esac

if [ "$VIGIL_VERSION" = "latest" ]; then
    VIGIL_VERSION=$(curl -s https://api.github.com/repos/sahil3982/vigil/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VIGIL_VERSION" ]; then
        echo "❌ Failed to fetch latest release version." >&2
        exit 1
    fi
fi

# Extract version number (remove 'v' prefix if present)
VERSION_NUMBER="${VIGIL_VERSION#v}"

# Try two possible URL patterns
URL1="https://github.com/sahil3982/vigil/releases/download/${VIGIL_VERSION}/vigil_${VERSION_NUMBER}_${OS}_${ARCH}.tar.gz"
URL2="https://github.com/sahil3982/vigil/releases/download/${VIGIL_VERSION}/vigil_v${VERSION_NUMBER}_${OS}_${ARCH}.tar.gz"

echo "📥 Downloading vigil ${VIGIL_VERSION} for ${OS}/${ARCH}..."

# Try first URL, fall back to second
tmpfile="$(mktemp -t vigil.XXXXXX.tar.gz)"
if curl -sfL "$URL1" -o "$tmpfile"; then
    echo "✅ Downloaded from pattern: vigil_${VERSION_NUMBER}_${OS}_${ARCH}.tar.gz"
elif curl -sfL "$URL2" -o "$tmpfile"; then
    echo "✅ Downloaded from pattern: vigil_v${VERSION_NUMBER}_${OS}_${ARCH}.tar.gz"
else
    echo "❌ Failed to download from either URL pattern" >&2
    echo "Tried:" >&2
    echo "  $URL1" >&2
    echo "  $URL2" >&2
    exit 1
fi

tar -xzf "$tmpfile" -C /tmp
rm -f "$tmpfile"

if [ "$OS" = "windows" ]; then
    BIN_NAME="vigil.exe"
    TARGET="$HOME/bin/$BIN_NAME"
    mkdir -p "$HOME/bin"
    install -m 755 "/tmp/$BIN_NAME" "$TARGET"
    echo "✅ Installed to: $TARGET"
    echo "💡 Add $HOME/bin to your PATH to run 'vigil' from anywhere."
else
    BIN_NAME="vigil"
    echo "🚚 Installing to $VIGIL_PREFIX..."
    sudo install -m 755 "/tmp/$BIN_NAME" "$VIGIL_PREFIX/$BIN_NAME"
fi

# Cleanup
rm -f "/tmp/vigil" "/tmp/vigil.exe"
echo "✅ Done! Try running: vigil --help"