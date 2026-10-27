---
type: log
owner: architect
surfaces: [cross]
feature: none
status: closed
opened: 2026-10-27
closed: 2026-10-27
supersedes:
verified-by: docs
---
# KB redesign from the CBAI sibling-KB audit

## Context
Audited the Orthox, GolfPlaying and Clinical-Photography KBs in `CBAI_Workspace` to find what to copy and what to avoid.

## Verified
- Entry files: Orthox `CLAUDE.md` 317 KB, GolfPlaying `CLAUDE.md` 360 KB; Clinical-Photography root `CLAUDE.md` belongs to a different project.
- Drift: React Native version stated three ways across docs in both Orthox and GolfPlaying; health scores disagree between README, registry and report in all three.
- Roll-up: GolfPlaying has 382 dated requirement docs, 17 about Wolf, with no single current-state page.
- Sync: no KB has a per-consumer status for API changes; admin (Orthox, GolfPlaying) or mobile (Clinical-Photography) is left out.
- Enforcement: no KB runs checks today. Clinical-Photography removed its CI; GolfPlaying's contract check lives in the mobile repo, outside the KB.

## Ask
- None. Kept ideas: Orthox structured API change notes (`registries/contract-changes.md`) and Verified/Ask requests (`templates/log-entry.md`); GolfPlaying Status/Supersedes headers and contract snapshot check (pending code); Clinical-Photography ledger, "Revisit when" on decisions, CI gate.

## Outcome
LL-005 to LL-008, D-005, `registries/owners.md`, `registries/contract-changes.md`, `features/`, `log/`, `templates/`, `scripts/kb_lint.py`.

## Updates
- 2026-10-27 (cross): created and closed.
