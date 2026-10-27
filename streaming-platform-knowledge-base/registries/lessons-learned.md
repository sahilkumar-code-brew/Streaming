---
type: registry
owner: architect
last-updated: 2026-10-27
---
# Lessons Learned (AUTHORITATIVE)

> Single source of truth for anti-patterns and their permanent fixes. Do NOT duplicate this table elsewhere. Reference it: `registries/lessons-learned.md`

| ID | Mistake | Consequence | Permanent fix | Enforced by |
|----|---------|-------------|---------------|-------------|
| LL-001 | Postgres queries in entitlement path | Pool exhaustion at 50k RPS | Redis O(1) key fetch, replicas for reads | `scripts/validate_code.sh` |
| LL-002 | WebSockets for score fan-out | Socket buffer bloat, reconnect storms | CDN-cached static JSON, 2s poll + jitter | `scripts/validate_code.sh` |
| LL-003 | Player/decoding logic on RN JS thread | Frame drops, UI freezes | Media3/AVPlayer in native modules, low-frequency events only | review |
| LL-004 | Fixed retry intervals | Thundering herd | Full jitter: `rand(0, min(cap, base * 2^attempt))` | review |

| LL-005 | Entry file (CLAUDE.md) used as changelog and status board | 300+ KB entry files in Orthox and GolfPlaying; agents cannot load them | CLAUDE.md is routing only, capped at 60 lines; history goes to `log/` and `contract-changes.md` | `scripts/kb_lint.py` |
| LL-006 | Same fact copied into several docs | Versions, health scores, milestones disagree across files in all three sibling KBs | Facts with IDs live in `registries/` only; other docs reference the ID | `scripts/kb_lint.py` |
| LL-007 | Dated fix docs never rolled up into current truth | 17 Wolf docs in GolfPlaying, no way to tell current state | Closing a log entry requires updating its `features/` page | `scripts/kb_lint.py` |
| LL-008 | API change shipped without per-consumer status | Admin and mobile silently miss changes | Every interface has a `contracts.md` row with a column per surface; changes logged in `contract-changes.md` | `scripts/check_kb_sync.sh` |

Add a row whenever the AI repeats a mistake. Next ID: LL-009.
