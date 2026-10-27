#!/bin/bash
set -euo pipefail

if ! git rev-parse --git-dir >/dev/null 2>&1; then
  echo "SKIP: not inside a git repository"
  exit 0
fi

comments=$(git diff --cached -U0 -- '*.ts' '*.tsx' '*.js' '*.go' '*.kt' '*.swift' \
  | grep '^+' | grep -v '^+++' \
  | grep -E '(^\+|[^:])//|/\*|^\+[[:space:]]*\*' \
  | grep -vE '^\+[[:space:]]*//(go:|[[:space:]]*\+build)' || true)

if [ -n "$comments" ]; then
  echo "ERROR: comments found; refactor into named helpers:"
  echo "$comments"
  exit 1
fi
echo "No comments in staged code."
