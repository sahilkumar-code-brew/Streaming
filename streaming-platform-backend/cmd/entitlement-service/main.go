package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/code-brew-ai/streaming-platform-backend/internal/handler/rest"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/config"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/database"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/health"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/keys"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/logger"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/ratelimit"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/server"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/session"
	redisrepo "github.com/code-brew-ai/streaming-platform-backend/internal/repository/redis"
	entitlementservice "github.com/code-brew-ai/streaming-platform-backend/internal/service/entitlement"
)

const (
	serviceName     = "entitlement-service"
	rateLimit       = 10
	rateLimitWindow = time.Minute
)

var version = "dev"

func main() {
	log := logger.New(serviceName, version)
	if err := run(log); err != nil {
		log.Error("service exited", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load[Config]()
	if err != nil {
		return err
	}
	ctx := context.Background()
	redisClient, err := database.NewRedis(ctx, cfg.RedisAddr)
	if err != nil {
		return err
	}
	defer redisClient.Close()

	sessionPublicKey, err := keys.Public(cfg.SessionPublicKey)
	if err != nil {
		return err
	}
	playbackPrivateKey, err := keys.PrivateFromSeed(cfg.PlaybackSigningKey)
	if err != nil {
		return err
	}

	repository := redisrepo.NewEntitlementRepository(redisClient)
	signer := entitlementservice.NewPlaybackSigner(playbackPrivateKey, cfg.PlaybackKeyID, cfg.PlaybackTokenTTL)
	settings := entitlementservice.Settings{
		OriginBaseURL:       cfg.OriginBaseURL,
		SteeringURL:         cfg.SteeringURL,
		TelemetrySampleRate: cfg.TelemetrySampleRate,
		RefreshLead:         cfg.PlaybackRefreshLead,
	}
	svc := entitlementservice.NewService(repository, signer, time.Now, settings)

	verifier := session.NewVerifier(sessionPublicKey)
	limiter := ratelimit.NewLimiter(rateLimit, rateLimitWindow, time.Now)
	readinessCheck := func(checkCtx context.Context) error { return redisClient.Ping(checkCtx).Err() }

	handler := rest.NewEntitlementRouter(svc, verifier, limiter, log, health.Check(readinessCheck))

	return server.Run(cfg.HTTPAddr, handler, log, cfg.ShutdownTimeout)
}
