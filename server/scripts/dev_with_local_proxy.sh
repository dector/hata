#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

go tool gust \
  -e 'HATA_DEV=1 go run -ldflags "-X hata/version.Version=$(./scripts/git_version.sh)" ./cmd/hata' \
  -TT \
  --e.before 'go tool templ generate ./internal/webui' \
  --exclude data \
  --exclude docs \
  --exclude out \
  --exclude .jj \
  --exclude.glob '*_templ.go'
