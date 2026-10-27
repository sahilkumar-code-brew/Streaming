---
type: feature
owner: architect
surfaces: [backend, mobile]
status: planned
last-verified: 2026-10-27
verified-by: docs
---
# Live Playback

> Current truth only. Per-surface build status lives in `registries/contracts.md`. Payloads live in `architecture/api-contracts.md`.

## Behaviour
A logged-in user (C-008) taps a live stream, gets a 300s playback token (C-001), and the native player starts on the CDN chosen by the steering manifest (C-002). The player renews the token before expiry without restarting playback. If the buffer drains, the player fails over to the next CDN in the priority list.

## Streaming targets
| Target | Value |
|--------|-------|
| Glass-to-glass latency | p95 <= 5s on LL-HLS; p95 <= 12s on DASH fallback |
| Segment duration | 2s CMAF fMP4, keyframe every 2s |
| LL-HLS parts | 0.5s, `PART-HOLD-BACK` 1.5s |
| Profiles | LL-HLS (iOS AVPlayer, Android Media3); DASH `urn:mpeg:dash:profile:cmaf:2019` from the same CMAF segments |
| Startup (tap to first frame) | p90 < 2s |
| Rebuffer ratio | < 0.5% of watch time |
| CDN failover trigger | buffer < 2 segments (4s) and a failed segment fetch, then next pathway |

## Telemetry policy (C-003)
| Setting | Value |
|---------|-------|
| Sampling | `telemetrySampleRate` from C-001, default 0.1; `error` events always sent |
| Beacon interval | flush every 30s |
| Batching window | 30s or 20 events, whichever comes first |
| Flush jitter | rand(0, 5s) added to every flush |
| Retry | 2 attempts, full jitter (LL-004), then drop the batch |
| Load at 20M viewers | 20M x 0.1 / 30s = ~67k req/s |

## Token lifecycle
1. `sessionJwt` (15 min) from C-008; refreshed at 80% of lifetime with jitter.
2. `playbackToken` (300s) from C-001; re-requested at `refreshAfterS + rand(0, 30)` s.
3. A failed renewal keeps playback on the old token until `exp + 30s`, retrying with full jitter.
4. origin-shield writes the token into segment URLs (D-007); the native module swaps tokens without a player restart.

## Surface notes
- **backend**: entitlement-service reads Redis `sub:{user_id}` only (LL-001); a miss returns 503 `SERVICE_UNAVAILABLE` (D-006). Signs Ed25519 (D-002). Edge or origin-shield verifies before serving.
- **mobile**: Media3/AVPlayer in native modules (LL-003); JS sees low-frequency buffer and stall events and passes only token strings down. Every retry uses full jitter (LL-004).
- **admin**: none directly; affected by steering weights (see `features/traffic-steering.md`).

## Open issues
- R-008: Ed25519 verification not confirmed on all three CDN edge runtimes.
- R-009: no concurrent device limit.
- R-010: Media3 HLS content steering support unverified.
- R-011: no fill path for Redis `sub:{user_id}`.

## History
- `log/2026-10-27-cross-org-standards-adoption.md`: org envelope, error codes, camelCase, UUIDs.
- `log/2026-10-27-cross-kb-contract-hardening.md`: schemas, token lifecycle, targets, telemetry policy.
- `log/2026-10-27-cross-kb-redesign.md`: page created from architecture docs.
