# Streaming Platform KB (learning project)

Practice KB for a Hotstar-scale live streaming platform with three surfaces: backend, mobile, admin. This file is routing only: no status, no changelog, no tables of facts (LL-005). Hard cap 60 lines, enforced by `scripts/kb_lint.py`.

## Read order
1. `registries/lessons-learned.md` and `instructions/code-standards.md` before generating code.
2. `features/<feature>.md` for current truth about whatever you are touching.
3. `registries/contracts.md` before writing code against any interface.

## Where does X go
| Content | Location | Rule |
|---------|----------|------|
| A fact with an ID or a version (stack, service, contract, decision, risk, lesson, owner) | `registries/` | One copy; everything else links by ID |
| Payload shapes | `architecture/api-contracts.md` | Every endpoint needs a `C-` row |
| History of interface changes | `registries/contract-changes.md` | Append-only, one row per change |
| Current behaviour of one feature | `features/<feature>.md` | Living page; edited, never appended to |
| Dated work: a fix, investigation, cross-team ask | `log/YYYY-MM-DD-<surface>-<topic>.md` | Closing it means updating its feature page |
| How it is built | `architecture/` | Max 300 lines per file |
| Standards, protocols, prompts | `instructions/` | |
| New doc skeletons | `templates/` | Copy, never link |

Surfaces are exactly `backend`, `mobile`, `admin`, `cross`. No `archive/` folder: superseded docs are deleted, git keeps history.

## Rules
- Cross-surface changes follow `instructions/change-protocol.md`.
- Every doc has frontmatter with `owner` (from `registries/owners.md`) and `last-verified`. Re-verify against code or live API, not other docs; record how in `verified-by`.
- Project rule: no code comments (overrides the org default of "comments for why").
- Run `python3 scripts/kb_lint.py` before finishing any KB edit. It must exit 0.
