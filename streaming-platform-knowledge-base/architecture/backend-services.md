---
type: architecture
owner: backend-lead
surfaces: [backend]
last-verified: 2026-10-26
verified-by: docs
---
# Backend Services

- Entitlement service: sub-3ms response, Ed25519 signing.
- Edge worker: validates token signature at the CDN edge before origin.
- Telemetry gateway: high-throughput Go service writing to Kafka.
