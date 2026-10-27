---
type: log
owner: architect
surfaces: [backend]
feature: features/name.md
status: open
opened: YYYY-MM-DD
closed:
supersedes:
verified-by: docs
---
# Short title

> File as `log/YYYY-MM-DD-<surface>-<topic>.md`, surface one of backend, mobile, admin, cross. `feature: none` only when no feature page applies. Closing: set `status: closed` and `closed:`, then update the feature page and bump its `last-verified` (kb_lint enforces this).

## Context
What was observed, by whom, where.

## Verified
Facts checked against code or live API, each with evidence (file:line, commit, request and response). Check the client side before blaming the backend.

## Ask
Who must do what. One bullet per surface.

## Outcome
Which feature page, registry row or contract change this produced.

## Updates
- YYYY-MM-DD (surface): dated note, appended as the work moves between teams.
