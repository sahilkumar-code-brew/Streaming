# Todo

- [ ] T1 — Restructure into domain/service/handler/repository/platform/pkg layout
- [ ] T2 — Add platform/database (redis), platform/ratelimit (+ apperror.CodeRateLimited), handler/rest/router.go
- [ ] T3 — cmd/entitlement-service/main.go wiring + .env.example
- [ ] T4 — go mod tidy (go.sum) + Makefile
- [ ] T5 — Table-driven tests (service, redis repo via testcontainers, validator); explicit D-006 503-on-miss test
- [ ] T6 — Design billing write-through + backfill job doc (closes R-011)
- [ ] T7 — KB updates: risks.md R-011, decisions.md D-014, features/live-playback.md, run kb_lint.py
