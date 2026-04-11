#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

go run ./cmd/devproxy -listen 127.0.0.1:4500 -target http://127.0.0.1:4501 &
proxy_pid=$!

cleanup() {
  kill "$proxy_pid" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

go tool \
  air \
    -build.cmd="go tool templ generate ./internal/webui && go build -o ./out/hata.air ./cmd/hata" \
    -build.full_bin="HATA_DEV=1 ./out/hata.air" \
    -build.include_dir="cmd,internal,pkg" \
    -build.include_ext="go,templ" \
    -build.exclude_regex="_templ.go" \
    -build.delay=100 \
    -tmp_dir=out
