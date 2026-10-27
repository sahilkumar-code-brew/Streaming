---
type: instruction
owner: architect
surfaces: [backend, mobile, admin]
last-verified: 2026-10-27
verified-by: docs
---
# Code Standards

1. No explanatory comments. Refactor into named helpers instead.
2. Names: camelCase vars/functions (verb-led), PascalCase types, SCREAMING_SNAKE_CASE constants/enums.
3. Max function length 25 lines, max nesting 2 (guard clauses).
4. No magic numbers or string literals; extract domain constants.
5. TypeScript/Go strict: no `any`, no unchecked assertions, explicit return types.
6. Layers: handler (parse, validate, delegate) -> service (logic) -> repository (queries, cache).
7. Typed custom errors only; catch blocks log context and re-throw or transform; never swallow.
8. No argument mutation; `const` by default.
9. Retries use full jitter (D-010); timeouts are explicit on every outbound call.

## Logging
- Structured JSON only. Every line carries `timestamp`, `level`, `message`, `requestId`, `traceId`, `spanId`, `service`, `version`.
- Levels: `error` needs attention, `warn` recovered, `info` business event, `debug` dev only.
- Never log passwords, OTPs, tokens (session, playback, refresh), DRM keys, full request bodies with PII, or connection strings. Log IDs instead.
- `X-Request-Id` is generated at the first hop if missing and propagated to every downstream call.

## Audit
Every admin action (C-005, C-006), every login and failed login, and every refresh-token reuse writes an append-only audit record: `timestamp`, `actorId`, `actorType`, `action`, `resourceType`, `resourceId`, `changes` (before/after), `ipAddress`.

## Services
- Every Go service exposes `GET /health/live` (process up) and `GET /health/ready` (Redis, Kafka, DB reachable).
- Handle SIGTERM: stop accepting, drain in-flight requests, close stores, then exit.
- Layout `cmd/`, `internal/`, `pkg/`; config from env via envconfig; never commit `.env`, only `.env.example`.
- Schema changes only through goose migrations; never edit a migration that has run.
- Kafka consumers follow D-012.

## Testing
- Go: table-driven tests with testify; stores tested with testcontainers, not mocks; `go test -race` in CI.
- Mobile: utils and hooks 90% line coverage, components 80%; player features also tested on a physical device.
- Admin: Vitest + RTL for components, Playwright for critical flows, axe-core for accessibility.

## Git
- Branches `feat/`, `fix/`, `chore/`, `refactor/`, `docs/`, `test/`, `hotfix/` + kebab-case; KB branches `docs/<user>/<topic>`.
- Conventional commits, imperative, subject under 72 chars; KB commits use `docs(kb): ...`.
- Code repo root holds only `README.md`, `CLAUDE.md`, optional `CHANGELOG.md`, `CONTRIBUTING.md`; docs go to this KB.

## Release blockers (org baseline automatic F)
Any of these blocks merge: secrets or `.env` in git history; hardcoded credentials; admin or protected route without auth; JWT verification disabled or `alg: none`; endpoint accepting input without validation; code repo without `README.md` and `CLAUDE.md`; production data without backups.
