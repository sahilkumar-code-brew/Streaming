---
type: feature
owner: architect
surfaces: [backend, mobile]
status: planned
last-verified: 2026-10-27
verified-by: docs
---
# Live Scores

> Current truth only. Per-surface build status lives in `registries/contracts.md`.

## Behaviour
The scorecard over the player updates every couple of seconds during a match.

## Contracts
- C-004: mobile polls a CDN-cached static JSON per stream (D-001, LL-002). Cache headers and poll interval are in `architecture/api-contracts.md`.

## Surface notes
- **backend**: scores-publisher owns the JSON (R-006, unbuilt).
- **mobile**: jittered polling hook; never WebSockets (LL-002).
- **admin**: none.

## Open issues
- R-006: scores-publisher must ship before mobile polls C-004.

## History
- `log/2026-10-27-cross-kb-redesign.md`: page created; missing producer flagged.
