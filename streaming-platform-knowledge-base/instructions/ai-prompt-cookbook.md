---
type: instruction
owner: architect
surfaces: [backend, mobile, admin]
last-verified: 2026-10-27
verified-by: docs
---
# AI Prompt Cookbook

Always prefix: "Read `registries/lessons-learned.md` and `instructions/code-standards.md` first. Build only against contracts in `architecture/api-contracts.md`."

1. Entitlement service (Go): "Using `architecture/api-contracts.md` section C-001 and `architecture/data-models.md`, build entitlement-service with `net/http` and go-redis: Redis-only subscription lookup (LL-001), 503 `SERVICE_UNAVAILABLE` on a Redis miss (D-006), Ed25519 playback token with the documented claims, org response envelope and error codes (D-008), `/health/live` and `/health/ready`, graceful SIGTERM shutdown."
2. Android player module: "Using `architecture/mobile-app.md` and `features/live-playback.md`, write the Media3 Expo Module (D-011) wrapped in `src/lib/`: buffer and stall events to JS at low frequency (LL-003), CDN failover on the trigger in the streaming targets table, steering priority from C-002, playback token swap without restart, full-jitter retries (D-010)."
3. Steering panel (Next.js): "Using `architecture/admin-dashboard.md` and `architecture/api-contracts.md` section C-005, build the Traffic Steering panel with sliders per `CdnPathway`, submitted through a TanStack Query mutation to C-005 with `If-Match` and `Idempotency-Key`. Validate integers summing to 100 before submit, ask for confirmation and send `confirmLargeShift` on a shift over 50 points, and handle `VERSION_CONFLICT` by refetching. Never write Redis directly (R-005)."
