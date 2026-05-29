#!/usr/bin/env bash
set -euo pipefail

REPO="navig-me/docker-doctor"
VERSION="${VERSION:-latest}"
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "unsupported arch: $ARCH"; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/deploy-doctor_${OS}_${ARCH}.tar.gz"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/deploy-doctor_${VERSION#v}_${OS}_${ARCH}.tar.gz"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
curl -fsSL "$URL" -o "$TMP_DIR/deploy-doctor.tar.gz"
tar -xzf "$TMP_DIR/deploy-doctor.tar.gz" -C "$TMP_DIR"
install -m 0755 "$TMP_DIR/deploy-doctor" /usr/local/bin/deploy-doctor

echo "deploy-doctor installed to /usr/local/bin/deploy-doctor"
