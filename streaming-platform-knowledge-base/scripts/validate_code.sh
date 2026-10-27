#!/bin/bash
set -euo pipefail

CODE_ROOT="${1:-.}"
STRICT="${STRICT:-0}"
violations=0
scanned=0

scan() {
  local dir="$CODE_ROOT/$1" id="$2"; shift 2
  if [ ! -d "$dir" ]; then
    echo "SKIP $id: $dir not found"
    return 0
  fi
  scanned=$((scanned + 1))
  if grep -rnE "$@" "$dir"; then
    echo "VIOLATION $id in $dir (see registries/lessons-learned.md)"
    violations=1
  fi
}

scan services/entitlement LL-001 -e 'SELECT ' -e 'FROM users' -e 'database/sql' -e 'pgx'
scan mobile/src LL-002 -e 'WebSocket' -e 'socket\.io'

if [ "$scanned" -eq 0 ] && [ "$STRICT" = "1" ]; then
  echo "FAIL: STRICT=1 and no code paths found under $CODE_ROOT"
  exit 1
fi
exit $violations
