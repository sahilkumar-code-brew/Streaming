---
type: instruction
owner: architect
surfaces: [backend, mobile, admin]
last-verified: 2026-10-27
verified-by: docs
---
# Cross-Surface Change Protocol

> Applies to any change that touches an interface in `../registries/contracts.md`. Contract first, code second.

## Order of work
1. **Registry**: add or update the row in `registries/contracts.md` (new ID, or bump version). Mark each affected surface `planned`.
2. **Shape**: update the payload in `architecture/api-contracts.md`. Every endpoint there must have a contract row.
3. **Change row**: append a row to `registries/contract-changes.md` (breaking, migration, commit, per-consumer status).
4. **Feature page**: update `features/<feature>.md` so it describes the new behaviour.
5. **Decision**: if the change alters a trade-off, add a row to `registries/decisions.md`.
6. **Backend ships first** and stays backward compatible with every consumer still on the previous version.
7. **Consumers** (mobile, admin) move to the new version; flip their column to `shipped` in both `contracts.md` and the `contract-changes.md` row.
8. **Deprecate** the old version only when every consumer column reads `shipped` on the new one. Mobile keeps the old version alive for at least one store-release cycle (old app builds stay in the wild).

## Breaking vs non-breaking
| Change | Breaking? | Action |
|--------|-----------|--------|
| Add optional response field | No | Same version |
| Add required request field, remove/rename field, change type or semantics | Yes | New version path (`/v2/...`), run both |
| Tighten validation, change error codes | Yes | Treat as new version |

## Definition of done (per contract change)
- [ ] `registries/contracts.md` row reflects real status on all three surfaces
- [ ] `architecture/api-contracts.md` shape matches the implementation
- [ ] `python3 scripts/kb_lint.py` exits 0 (it runs `scripts/check_kb_sync.sh` too)
- [ ] No surface reads or writes a store owned by another service (admin goes through backend APIs, except read-only analytics like C-007)
- [ ] Lessons check: `registries/lessons-learned.md` reviewed

## When code exists
Snapshot the live API into `api-contract.snapshot.json` and fail mobile and admin CI on drift (GolfPlaying's `api_contract_check.py` pattern, but run from CI, not by hand).
