---
type: registry
owner: backend-lead
last-updated: 2026-10-27
---
# Contract Changes (AUTHORITATIVE, append-only)

> One row per change to an interface in `registries/contracts.md`. Never put change banners at the top of `architecture/api-contracts.md` (LL-005). Reference it: `registries/contract-changes.md`

Breaking: `yes` requires a new version path and both versions running (see `instructions/change-protocol.md`). Commit: backend commit that implements it, or `unbuilt`. Migration: DB/Redis migration needed, or `none`.

| ID | Date | Contract | Change | Breaking | Migration | Commit | Mobile | Admin |
|----|------|----------|--------|----------|-----------|--------|--------|-------|
| CC-001 | 2026-10-26 | C-001 | Initial v1: entitlement returns token, manifest, steering URL | no | none | unbuilt | planned | n/a |
| CC-002 | 2026-10-27 | C-005 | Proposed: steering weights write moves from direct Redis to steering-service API | no | none | unbuilt | n/a | planned |
| CC-003 | 2026-10-27 | C-001 | Full v1 shape, problem+json errors, playback token refresh | no | none | unbuilt | planned | n/a |
| CC-004 | 2026-10-27 | C-002, C-004 | Paths moved under `/api/v1/`; C-004 owner scores-publisher | no | none | unbuilt | planned | n/a |
| CC-005 | 2026-10-27 | C-005, C-006 | Paths, payloads, optimistic concurrency | no | none | unbuilt | n/a | planned |
| CC-006 | 2026-10-27 | C-008, C-009 | Initial v1: OTP login and session token | no | none | unbuilt | planned | planned |
| CC-007 | 2026-10-27 | C-001, C-003, C-004, C-005, C-006, C-008, C-009 | Org envelope and error codes, camelCase fields, UUID v4 IDs (D-008) | no | none | unbuilt | planned | planned |

Next ID: CC-008.
