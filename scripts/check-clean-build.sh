#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

printf '%s\n' 'Building the web frontend and generating src/dist...'
cd "$ROOT_DIR/web"
npm ci --include=dev
npm run build

printf '%s\n' 'Running Go tests against the generated embed assets...'
cd "$ROOT_DIR/src"
go test ./...
go build ./...

printf '%s\n' 'Clean checkout verification passed.'
