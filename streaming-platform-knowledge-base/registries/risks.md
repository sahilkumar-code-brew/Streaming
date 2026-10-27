---
type: registry
owner: architect
last-updated: 2026-10-27
---
# Risks (AUTHORITATIVE)

> Single source of truth for risks. Reference it: `registries/risks.md`

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|-----------|--------|------------|
| R-001 | Single CDN outage during a match | Med | High | Multi-CDN with content steering + client failover |
| R-002 | Reconnect/retry storm after an incident | Med | High | Full-jitter backoff everywhere (LL-004) |
| R-003 | Redis hot-key on popular stream metadata | Med | Med | Local in-process cache with short TTL, key sharding |
| R-004 | Signing key compromise | Low | High | Key rotation with `kid`, short token TTL (300s) |
| R-005 | Admin writes Redis `steering:weights` directly | High | High | Route through steering-service API (C-005): validation, sum-to-100 check, audit log |
| R-006 | Live scores producer unbuilt | Med | Med | scores-publisher owns C-004; build before mobile polls it |
| R-007 | Stream ops (start/stop, DRM) unbuilt | Med | High | stream-ops-service owns C-006; build before admin UI |
| R-008 | Ed25519 verify not available on every CDN edge runtime | Med | High | Verify per vendor (Akamai, CloudFront, Fastly) before building edge-worker; fallback is verify at origin-shield |
| R-009 | No concurrent device limit per subscription | High | High | Business decision pending; `ENTITLEMENT_DEVICE_LIMIT` code reserved |
| R-010 | Media3 support for HLS content steering unverified | Med | Med | Prototype in Phase 1; if missing, native module applies C-002 priority itself |
| R-011 | No defined fill path for Redis `sub:{user_id}` from billing | High | High | Define billing write-through plus backfill job before entitlement ships (D-006) |
