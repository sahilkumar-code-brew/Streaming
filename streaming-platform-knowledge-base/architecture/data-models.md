---
type: architecture
owner: backend-lead
surfaces: [backend]
last-verified: 2026-10-26
verified-by: docs
---
# Data Models

- Postgres (billing only): users, subscriptions, streams, cdn_configurations.
- Redis keys: `sub:{user_id}`, `stream:meta:{stream_id}`, `steering:weights`.
- ClickHouse: `client_qoe_telemetry`, ordered by (cdn_provider, stream_id, event_timestamp).
