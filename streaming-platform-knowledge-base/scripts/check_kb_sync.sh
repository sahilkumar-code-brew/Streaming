#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."

CONTRACTS=registries/contracts.md
SERVICES=registries/services.md
API_DOC=architecture/api-contracts.md
failures=0

fail() {
  echo "FAIL: $1"
  failures=1
}

check_documented_endpoints_have_contracts() {
  local path
  while read -r path; do
    grep -qF " $path\`" "$CONTRACTS" || fail "$API_DOC endpoint $path has no row in $CONTRACTS"
  done < <(grep -oE '`(GET|POST|PUT|PATCH|DELETE) [^`?]+' "$API_DOC" | sed -E 's/^`[A-Z]+ //' | sed -E 's#^https?://[^/]+##')
}

check_contracts_have_schemas() {
  local id endpoint
  while IFS='|' read -r id endpoint; do
    grep -qE "^## $id " "$API_DOC" || fail "$CONTRACTS row $id has no '## $id' section in $API_DOC"
    [ -z "$endpoint" ] && continue
    grep -qF "\`$endpoint\`" "$API_DOC" || grep -qF "\`$endpoint?" "$API_DOC" || fail "$CONTRACTS row $id endpoint '$endpoint' not documented in $API_DOC"
  done < <(awk -F'|' '/^\| C-[0-9]+/ {gsub(/ /,"",$2); ep=$3; if (match(ep, /`(GET|POST|PUT|PATCH|DELETE) [^`]+`/)) ep=substr(ep, RSTART+1, RLENGTH-2); else ep=""; print $2 "|" ep}' "$CONTRACTS")
}

check_schemas_have_contracts() {
  local id
  while read -r id; do
    grep -qE "^\| $id \|" "$CONTRACTS" || fail "$API_DOC section $id has no row in $CONTRACTS"
  done < <(grep -oE '^## C-[0-9]+' "$API_DOC" | sed 's/^## //')
}

check_owners_are_registered_services() {
  local owner
  while read -r owner; do
    grep -qE "^\| $owner \|" "$SERVICES" || fail "contract owner '$owner' missing from $SERVICES"
  done < <(awk -F'|' '/^\| C-[0-9]+/ {gsub(/ /,"",$5); print $5}' "$CONTRACTS" | grep -v '\*\*unowned\*\*' | sort -u)
}

check_unowned_contracts() {
  local unowned
  unowned=$(grep -E '^\| C-[0-9]+' "$CONTRACTS" | grep -F '**unowned**' | awk -F'|' '{print $2}' | tr -d ' ' | tr '\n' ' ' || true)
  [ -z "$unowned" ] || echo "WARN: unowned contracts: $unowned"
}

check_documented_endpoints_have_contracts
check_contracts_have_schemas
check_schemas_have_contracts
check_owners_are_registered_services
check_unowned_contracts

[ "$failures" -eq 0 ] && echo "KB sync OK."
exit $failures
