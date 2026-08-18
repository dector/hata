#!/usr/bin/env bash
set -euo pipefail

# Print an application version derived from git tags.
#
# Output examples:
#   v1.2.3                  exact tag
#   v1.2.3-4-gabcdef0       4 commits after tag
#   v1.2.3-4-gabcdef0-dirty uncommitted changes
#   dev                     no git tag available

cd "$(dirname "$0")/.."

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  printf 'dev\n'
  exit 0
fi

is_dirty() {
  ! git diff --quiet --ignore-submodules -- 2>/dev/null || \
    ! git diff --cached --quiet --ignore-submodules -- 2>/dev/null || \
    [ -n "$(git ls-files --others --exclude-standard)" ]
}

if ! git describe --tags --abbrev=0 >/dev/null 2>&1; then
  version="v0.0.0-$(git rev-parse --short HEAD)"
  if is_dirty; then
    version="${version}-dirty"
  fi
  printf '%s\n' "${version}"
  exit 0
fi

version="$(git describe --tags --always)"
if is_dirty; then
  version="${version}-dirty"
fi
printf '%s\n' "${version}"
