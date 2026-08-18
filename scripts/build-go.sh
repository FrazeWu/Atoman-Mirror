#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build"

mkdir -p "$BUILD_DIR"

printf '%s\n' 'Building the web frontend...'
cd "$ROOT_DIR/web"
npm ci --include=dev
npm run build

printf '%s\n' 'Building the Go binary...'
cd "$ROOT_DIR/src"
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$BUILD_DIR/hubproxy" .

printf 'Built %s\n' "$BUILD_DIR/hubproxy"
