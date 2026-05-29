#!/usr/bin/env bash
set -euo pipefail

REPO="navig-me/deploy-doctor"
VERSION="${VERSION:-latest}"
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "unsupported arch: $ARCH"; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  LATEST_TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
  if [ -z "$LATEST_TAG" ]; then
    echo "failed to resolve latest release tag"
    exit 1
  fi
  VERSION_NUM="${LATEST_TAG#v}"
  URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/deploy-doctor_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
else
  VERSION_NUM="${VERSION#v}"
  URL="https://github.com/${REPO}/releases/download/${VERSION}/deploy-doctor_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
curl -fsSL "$URL" -o "$TMP_DIR/deploy-doctor.tar.gz"
tar -xzf "$TMP_DIR/deploy-doctor.tar.gz" -C "$TMP_DIR"
install -m 0755 "$TMP_DIR/deploy-doctor" /usr/local/bin/deploy-doctor

echo "deploy-doctor installed to /usr/local/bin/deploy-doctor"
