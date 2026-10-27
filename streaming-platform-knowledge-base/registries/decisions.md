---
type: registry
owner: architect
last-updated: 2026-10-27
---
# Decisions (AUTHORITATIVE)

> Single source of truth for decisions. Reference it: `registries/decisions.md`. A row that deviates from the org CB standards (core/conventions.md, profiles/) is that deviation's ADR and says so in Why.

| ID | Date | Decision | Why | Revisit when |
|----|------|----------|-----|--------------|
| D-001 | 2026-10-26 | Scores via CDN-cached JSON polling, not WebSockets | 20M sockets are costly; cacheable at edge (LL-002) | Scores need sub-second latency or per-user data |
| D-002 | 2026-10-26 | Ed25519 tokens verified at the edge | Asymmetric: edge needs only the public key | Edge platform cannot verify Ed25519 |
| D-003 | 2026-10-26 | Postgres off the playback path | LL-001 | Never |
| D-004 | 2026-10-26 | Zero explanatory comments in code | Self-documenting names; see `instructions/code-standards.md` | Team agrees otherwise |
| D-005 | 2026-10-27 | KB is enforced by `scripts/kb_lint.py`, not by convention | Sibling CBAI KBs drifted wherever rules were unenforced (LL-005..LL-008) | Lint becomes slower than the value it adds |
| D-006 | 2026-10-27 | Redis miss on `sub:{user_id}` returns 503 `SERVICE_UNAVAILABLE` (retryable) and alerts | Never paywall a paying user because cache data is missing; Postgres fallback forbidden (LL-001) | A verified fill path makes a miss mean "not subscribed" |
| D-007 | 2026-10-27 | Playback token carried in segment URLs via manifest rewrite at origin-shield | Player-agnostic; AVPlayer has no documented per-request header API | All target players support documented custom headers |
| D-008 | 2026-10-27 | All JSON APIs use the org envelope `{success, data, error, meta}`, org error codes plus domain codes, camelCase fields, UUID v4 IDs, ISO 8601 UTC. Exceptions: C-002 (HLS steering spec), C-004 (static CDN file) | Org CB conventions; corrected the same day from RFC 9457, which renamed the org codes | Org conventions change |
| D-009 | 2026-10-27 | Every HTTP contract path lives under `/api/v1/` (or `/api/vN/`) | One routing and versioning rule for Envoy and the CDN | Never |
| D-010 | 2026-10-27 | Every retry uses full jitter (LL-004), not the org fixed 1s/2s/4s backoff | Deviation from org conventions: fixed intervals synchronise 20M clients into retry storms (R-002) | Never |
| D-011 | 2026-10-27 | Mobile is Expo with prebuild and dev client; Media3/AVPlayer players are Expo Modules (Kotlin/Swift) wrapped in `src/lib/` | Org react-native profile default is Expo; Expo Modules API covers custom native players, so no bare-RN deviation | A native requirement Expo Modules cannot meet |
| D-012 | 2026-10-27 | Kafka events use CloudEvents 1.0, names `{domain}.{entity}.{past-tense-verb}`, idempotent consumers, DLQ with alert and no auto-replay | Org microservices profile | Never |
| D-013 | 2026-10-27 | Backend services stay in Go; Rust only for a Fastly Compute edge-worker (Wasm), if built | Services are I/O-bound (Redis, Kafka, HTTP), so Go meets p99 targets with faster delivery and an org go-service profile; Rust is first-class on Fastly Compute (R-008) | A service becomes CPU-bound (in-house packaging or media parsing) or GC pauses break a p99 target |

Next ID: D-014.
