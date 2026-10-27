---
type: feature
owner: architect
surfaces: [backend, mobile, admin]
status: planned
last-verified: 2026-10-27
verified-by: docs
---
# Traffic Steering

> Current truth only. Per-surface build status lives in `registries/contracts.md`.

## Behaviour
Operators split live traffic across Akamai, CloudFront and Fastly (split in `registries/tech-stack.md`). Players pick up the new priority on their next steering manifest fetch.

## Contracts
- C-002: players read the priority list.
- C-005: admin writes weights (proposed; does not exist yet).

## Surface notes
- **backend**: steering-service builds the RFC 8216bis manifest from Redis `steering:weights`. Should validate weights sum to 100 and audit every write.
- **mobile**: re-fetches the steering manifest and reorders CDN failover.
- **admin**: sliders per CDN. Today the design writes Redis directly, bypassing validation and audit (R-005).

## Open issues
- R-005: admin must move to the C-005 API before admin builds the panel.

## History
- `log/2026-10-27-cross-kb-redesign.md`: page created; direct Redis write flagged.
