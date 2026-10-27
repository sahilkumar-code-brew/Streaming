---
type: architecture
owner: architect
surfaces: [backend, mobile, admin]
last-verified: 2026-10-26
verified-by: docs
---
# Core Architecture

> Stack choices: `../registries/tech-stack.md`. Services: `../registries/services.md`.

- Targets: 20M concurrent streams, 50 Tbps egress.
- Control plane: Go entitlement APIs behind Envoy, Redis Cluster.
- Data plane: SRT ingest -> CMAF fMP4 transcoding -> Varnish origin shield -> multi-CDN with content steering.
- Client: React Native UI, native Media3/AVPlayer.
- Admin: Next.js with ClickHouse metrics.

## Invariants
- No relational DB call during token issuance.
- No WebSockets for scores.
- Every network call uses jittered exponential backoff.
- Controllers thin; logic in services or Zustand stores.
