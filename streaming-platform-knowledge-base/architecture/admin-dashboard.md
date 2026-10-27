---
type: architecture
owner: admin-lead
surfaces: [admin]
last-verified: 2026-10-27
verified-by: docs
---
# Admin Web Dashboard

- Next.js App Router, Tailwind, TanStack Query.
- Stream operations: state changes through C-006 (owned by stream-ops-service).
- Traffic steering: sliders calling the C-005 API; never Redis directly (R-005).
- NOC: live ClickHouse queries for per-CDN rebuffer percentage.
