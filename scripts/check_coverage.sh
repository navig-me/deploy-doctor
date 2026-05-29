#!/usr/bin/env bash
set -euo pipefail

THRESHOLD="${1:-60}"
TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

go test ./... -coverprofile="$TMP" >/dev/null
COV="$(go tool cover -func="$TMP" | awk '/^total:/ {gsub("%", "", $3); print $3}')"

awk -v cov="$COV" -v th="$THRESHOLD" 'BEGIN { if (cov+0 < th+0) exit 1 }' || {
  echo "coverage gate failed: ${COV}% < ${THRESHOLD}%"
  exit 1
}

echo "coverage gate passed: ${COV}% >= ${THRESHOLD}%"
