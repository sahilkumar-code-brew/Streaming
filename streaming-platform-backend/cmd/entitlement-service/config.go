package main

import "time"

type Config struct {
	HTTPAddr            string        `envconfig:"HTTP_ADDR" default:":8081"`
	RedisAddr           string        `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	SessionPublicKey    string        `envconfig:"SESSION_PUBLIC_KEY" required:"true"`
	PlaybackSigningKey  string        `envconfig:"PLAYBACK_SIGNING_KEY" required:"true"`
	PlaybackKeyID       string        `envconfig:"PLAYBACK_KEY_ID" required:"true"`
	PlaybackTokenTTL    time.Duration `envconfig:"PLAYBACK_TOKEN_TTL" default:"300s"`
	PlaybackRefreshLead time.Duration `envconfig:"PLAYBACK_REFRESH_LEAD" default:"60s"`
	OriginBaseURL       string        `envconfig:"ORIGIN_BASE_URL" required:"true"`
	SteeringURL         string        `envconfig:"STEERING_URL" required:"true"`
	TelemetrySampleRate float64       `envconfig:"TELEMETRY_SAMPLE_RATE" default:"0.1"`
	ShutdownTimeout     time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"15s"`
}
