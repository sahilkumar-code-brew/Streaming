---
type: registry
owner: backend-lead
last-updated: 2026-10-27
---
# Services (AUTHORITATIVE)

> Single source of truth for the service roster. Reference it: `registries/services.md`

| Service | Language | Responsibility | Store | Target | Status | Owner |
|---------|----------|----------------|-------|--------|--------|-------|
| entitlement-service | Go | Validate subscription, mint Ed25519 token | Redis | p99 < 3ms | planned | backend-lead |
| edge-worker | CDN worker | Verify token signature at edge | none | | planned | backend-lead |
| telemetry-gateway | Go | Ingest QoE beacons into Kafka | Kafka | | planned | backend-lead |
| steering-service | Go | Serve steering manifest from weights | Redis | | planned | backend-lead |
| admin-web | Next.js | Stream ops, steering sliders, NOC view | ClickHouse, Redis | | planned | admin-lead |
| mobile-app | React Native | Playback, scores, native telemetry | | | planned | mobile-lead |
| auth-service | Go | OTP login, issue and rotate session JWT (C-008, C-009); RBAC roles `viewer`, `ops`, `steering-admin` in `roles` claim | Postgres (users), Redis (refresh tokens) | p99 < 50ms | planned | backend-lead |
| catalog-service | Go | Stream list and metadata, writes `stream:meta:{stream_id}` | Postgres, Redis | | planned | backend-lead |
| stream-ops-service | Go | Sole writer of stream state (C-006) | Postgres, Redis | | planned | backend-lead |
| ingest-packager | FFmpeg + packager | SRT ingest, ABR encode, CMAF LL-HLS and DASH packaging | object storage | | planned | backend-lead |
| origin-shield | Varnish | Origin shield in front of packager output, manifest token rewrite (D-007) | none | | planned | backend-lead |
| telemetry-consumer | Go | Kafka -> ClickHouse batch insert | ClickHouse | | planned | backend-lead |
| scores-publisher | Go | Publishes per-stream score JSON to origin (C-004) | object storage | | planned | backend-lead |
| drm-license-proxy | Go | Widevine/FairPlay license proxy, checks playback token | Redis | | planned | backend-lead |
