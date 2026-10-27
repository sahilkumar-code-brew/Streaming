# Plan: entitlement-service runnable + billing fill path (closes R-011)

## Context
- Only `internal/entitlement/*.go` exists (handler/service/token/model/repository all in one
  flat package). It correctly implements C-001 but has no `cmd/` entrypoint, no redis client
  wiring, no rate limiting, and no `go.sum`.
- 12 other services in `registries/services.md` are `planned`, no code.
- R-011 (`registries/risks.md`): no fill path for Redis `sub:{user_id}` from billing.
  D-006 requires this be defined before entitlement-service ships.
- User supplied a target package layout (domain/service/handler/repository/platform/pkg,
  `cmd/api/main.go`, `go.sum`, `Makefile`, `migrations/`). This plan adopts that layout,
  adapted:
  - `cmd/entitlement-service/main.go` instead of `cmd/api/main.go` — matches D-013 / the
    service registry (independent Go binaries per service, not one monolith).
  - No `migrations/` yet — entitlement-service is Redis-only (LL-001: Postgres off the
    playback path). Migrations land with the first Postgres-backed service (auth-service
    or whatever owns billing).
  - `pkg/validator` holds the stream_id/device_id format checks currently inline in the
    handler.

## Target layout
```
streaming-platform-backend/
├── cmd/
│   └── entitlement-service/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── entitlement/
│   │       ├── entity.go        # Subscription, Stream, Snapshot, Grant, Request (was model.go)
│   │       └── repository.go    # Repository interface
│   ├── service/
│   │   └── entitlement/
│   │       └── service.go       # Grant logic (was entitlement/service.go)
│   ├── handler/
│   │   └── rest/
│   │       ├── entitlement_handler.go
│   │       └── router.go        # new: mux + middleware chain + health, extracted from main
│   ├── repository/
│   │   └── redis/
│   │       └── entitlement_repo.go  # was entitlement/repository.go
│   └── platform/
│       ├── config/
│       │   └── config.go        # was entitlement/config.go, generalized
│       ├── database/
│       │   └── redis.go         # new: redis.Client constructor (did not exist)
│       ├── logger/
│       │   └── logger.go        # renamed from platform/logging
│       ├── ratelimit/
│       │   └── ratelimit.go     # new: 10/min per user, C-001 (gap found in review)
│       ├── httpx/                # unchanged
│       ├── apperror/              # unchanged
│       ├── session/               # unchanged
│       ├── keys/                  # unchanged
│       └── server/                # unchanged
├── pkg/
│   └── validator/
│       └── validator.go         # new: stream_id pattern, device_id UUIDv4 check
├── go.mod
├── go.sum                        # new: go mod tidy
└── Makefile                      # new
```

`internal/entitlement/token.go` moves to `internal/service/entitlement/token.go` (it's
signing logic the service layer owns, not a separate domain concern).

## Tasks

### T1 — Restructure entitlement code into the new layout
Move/rename, no behavior change:
- `model.go` → `internal/domain/entitlement/entity.go`
- `repository.go` interface → `internal/domain/entitlement/repository.go`;
  `RedisRepository` impl → `internal/repository/redis/entitlement_repo.go`
- `service.go`, `token.go` → `internal/service/entitlement/`
- `handler.go` → `internal/handler/rest/entitlement_handler.go`
- `config.go` → `internal/platform/config/config.go`
- `internal/platform/logging/` → `internal/platform/logger/` (package renamed `logger`)
- Extract `streamIDPattern` regex + device-id UUIDv4 check from the handler into
  `pkg/validator/validator.go` (`ValidateStreamID`, `ValidateDeviceID`)
- Fix import paths everywhere they broke.

**Acceptance:** `go build ./...` succeeds with the new paths; no package left in the old
`internal/entitlement/` location; `go vet ./...` clean.

### T2 — Add the missing platform pieces
- `internal/platform/database/redis.go`: `New(addr string) (*redis.Client, error)` with an
  explicit `PING` on startup (fail fast) and a `Close` used on SIGTERM drain.
- `internal/platform/ratelimit/ratelimit.go`: per-user token bucket or fixed-window
  limiter, 10 req/min (C-001), keyed on `session.UserID(ctx)`; on limit hit returns
  `apperror.New(apperror.CodeRateLimited, ...)` with `retryAfterMs` — **note:**
  `apperror.CodeRateLimited` doesn't exist yet either; add it (`RATE_LIMITED`, 429,
  retryable) alongside the limiter.
- `internal/handler/rest/router.go`: builds the `*http.ServeMux`, registers
  `entitlement_handler` under `entitlement.RoutePattern`, wraps with
  `session.Require` + rate-limit middleware + `httpx.Chain(WithRequestID, WithAccessLog,
  WithRecovery)`, calls `health.Register`.

**Acceptance:** unit test proves the 11th request in a rolling minute for one user gets
`429 RATE_LIMITED` with `retryAfterMs > 0`; a different user's 1st request in the same
window succeeds. `go test ./internal/platform/ratelimit/...` passes.

### T3 — cmd/entitlement-service/main.go
Wiring only (no new logic): `config.Load` → `database.NewRedis` → `keys.PrivateFromSeed` /
`keys.Public` (session verify key) → build `Repository`, `TokenSigner`, `Service` →
`router.New(...)` → `server.Run(cfg.HTTPAddr, mux, logger, cfg.ShutdownTimeout)`. Logger via
`logger.New("entitlement-service", version)` where `version` is a build-time `ldflags` var
(not yet defined anywhere — add `var version = "dev"`).

**Acceptance:** `go run ./cmd/entitlement-service` starts and serves `GET /health/live` ->
200 with real env vars (or documented required envs in a new `.env.example`, never
committing `.env` per code-standards.md). `go build ./...` produces a binary.

### T4 — go.sum + Makefile
- `go mod tidy` to generate `go.sum` (currently missing — nothing in `go.mod` has a
  `require` block even though 4 third-party modules are imported).
- `Makefile` targets: `build`, `run`, `test` (`go test -race ./...` per
  instructions/code-standards.md Testing section), `lint` (placeholder if no linter
  configured yet), `tidy`.

**Acceptance:** `make test` runs `go test -race ./...` and passes.

### T5 — Table-driven tests for what's being touched
Per code-standards.md ("Go: table-driven tests with testify; stores tested with
testcontainers, not mocks"):
- `internal/service/entitlement/service_test.go`: table-driven over
  `checkEligibility` branches (not found, not live, no subscription row ⇒ 503 per D-006,
  expired subscription, happy path).
- `internal/repository/redis/entitlement_repo_test.go`: testcontainers-redis, not mocks.
- `pkg/validator/validator_test.go`: valid/invalid stream_id and device_id cases.

**Acceptance:** `go test -race ./...` green; the D-006 "Redis miss returns 503, not a
paywall" behavior has an explicit test (it's the contract's single most important
invariant and isn't tested today).

### T6 — Design the billing fill path (closes R-011)
This is the actual blocker for D-006/entitlement shipping, not just scaffolding. Minimal
design for a learning project:
- **Write-through:** billing (stubbed — no billing-service exists yet) writes
  `sub:{user_id}` to Redis synchronously on every subscription state change
  (`status`, `expires_at` fields, matching `repository.go`'s `toSubscription`).
  Document as a new contract or an internal event
  (`billing.subscription.updated`, CloudEvents per D-012) that a small consumer applies
  to Redis — pick whichever is simpler to actually stub; a direct write-through function
  is enough for this learning project, call it out as a simplification vs. production.
- **Backfill job:** one-shot script/cmd (`cmd/billing-backfill/main.go` or a documented
  manual SQL→Redis dump) that seeds `sub:{user_id}` for all currently-active subscriptions
  before entitlement-service's first real traffic, so a Redis miss only ever means
  "truly not subscribed," never "never backfilled."
- Write this up as `architecture/billing-fill-path.md` (new, backend surface) or fold into
  an existing architecture doc if `kb_lint.py`'s structure prefers that — check lint rules
  first.

**Acceptance:** doc exists, passes `python3 scripts/kb_lint.py`; a concrete fill mechanism
is named (not "TBD").

### T7 — Update the KB per instructions/change-protocol.md
- `registries/risks.md`: R-011 mitigation updated from "Define... before entitlement
  ships (D-006)" to reference the new doc; likelihood/impact re-assessed once a design
  exists (still "planned," not built — don't mark resolved, there's no code yet).
- `registries/decisions.md`: new `D-014` recording the write-through + backfill choice,
  `Why`, `Revisit when`.
- `features/live-playback.md`: "Open issues" R-011 line gets a pointer to the new doc;
  "History" gets a new bullet.
- `registries/contracts.md` / `registries/services.md`: if the write-through becomes a
  formal contract (`C-010`?) or the backfill job becomes a named service/cmd, add rows —
  otherwise leave as an internal implementation detail of entitlement-service.
- Run `python3 scripts/kb_lint.py` and (if it exists and is relevant)
  `scripts/check_kb_sync.sh` before calling this closed.

**Acceptance:** both scripts exit 0; `features/live-playback.md` no longer lists R-011 as
an open, undefined issue.

## Dependency order
T1 → T2 → T3 → T4 → T5 (code path, each step buildable/testable before the next)
T6 → T7 (KB path, can start in parallel with T1 once T1's restructuring doesn't block it)

## Checkpoints
- **After T3:** entitlement-service actually runs end-to-end locally against a real Redis
  — manually verify a 200 from C-001 and a 503 from a missing `sub:{user_id}` key (D-006)
  before writing tests in T5.
- **After T5:** `go test -race ./...` and `make test` both green before touching the KB.
- **After T7:** `kb_lint.py` exit 0 before considering this plan done.
