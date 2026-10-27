#!/usr/bin/env bash
set -euo pipefail

if ! node -e 'process.exit(Number(process.versions.node.split(".")[0]) >= 24 ? 0 : 1)'; then
  echo 'Node.js 24 or newer is required; run nvm use before building.' >&2
  exit 1
fi

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
