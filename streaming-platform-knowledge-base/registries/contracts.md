---
type: registry
owner: architect
last-updated: 2026-10-27
---
# Contracts (AUTHORITATIVE)

> Single source of truth for which surface owns and consumes each interface. Payload shapes live in `architecture/api-contracts.md`; this table tracks ownership and per-surface sync state. Reference it: `registries/contracts.md`

Status values: `planned` -> `in-progress` -> `shipped` -> `deprecated` -> `removed`. `n/a` means the surface does not consume it.

| ID | Interface | Version | Owner (producer) | Mobile | Admin | Backend | Notes |
|----|-----------|---------|------------------|--------|-------|---------|-------|
| C-001 | `POST /api/v1/streams/{stream_id}/entitlement` | v1 | entitlement-service | planned | n/a | planned | Redis only (LL-001) |
| C-002 | `GET /api/v1/steering.json` | v1 | steering-service | planned | n/a | planned | Player reads; admin writes weights via C-005 |
| C-003 | `POST /api/v1/telemetry/qoe` | v1 | telemetry-gateway | planned | n/a | planned | |
| C-004 | `GET /api/v1/scores/{stream_id}.json` | v1 | scores-publisher | planned | n/a | planned | CDN static (R-006) |
| C-005 | `PUT /api/v1/admin/steering/weights` | v1 | steering-service | n/a | planned | planned | Replaces direct Redis write (R-005) |
| C-006 | `PUT /api/v1/admin/streams/{stream_id}/state` | v1 | stream-ops-service | n/a | planned | planned | R-007 |
| C-007 | NOC rebuffer % (ClickHouse read) | n/a | admin-web | n/a | planned | n/a | Direct read-only query; acceptable |
| C-008 | `POST /api/v1/auth/token` | v1 | auth-service | planned | planned | planned | Issues and rotates session JWT |
| C-009 | `POST /api/v1/auth/otp` | v1 | auth-service | planned | planned | planned | Sends login OTP |

Add a row before any surface writes code against a new interface. Next ID: C-010.
