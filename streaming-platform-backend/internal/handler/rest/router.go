package rest

import (
	"log/slog"
	"net/http"

	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/health"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/httpx"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/ratelimit"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/session"
)

func NewEntitlementRouter(granter Granter, verifier *session.Verifier, limiter *ratelimit.Limiter, logger *slog.Logger, healthChecks ...health.Check) http.Handler {
	mux := http.NewServeMux()
	entitlementHandler := NewEntitlementHandler(granter, logger)
	protected := httpx.Chain(entitlementHandler, session.Require(verifier, logger), ratelimit.Middleware(limiter, logger))
	mux.Handle(EntitlementRoutePattern, protected)
	health.Register(mux, healthChecks...)
	return httpx.Chain(mux, httpx.WithRequestID, httpx.WithRecovery(logger), httpx.WithAccessLog(logger))
}
