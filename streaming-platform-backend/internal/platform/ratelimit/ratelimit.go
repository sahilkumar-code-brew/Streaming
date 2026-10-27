package ratelimit

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/apperror"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/httpx"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/session"
)

type window struct {
	startedAt time.Time
	count     int
}

type Limiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	now     func() time.Time
	windows map[string]*window
}

func NewLimiter(limit int, interval time.Duration, now func() time.Time) *Limiter {
	return &Limiter{limit: limit, window: interval, now: now, windows: make(map[string]*window)}
}

func (l *Limiter) Allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	current, exists := l.windows[key]
	if !exists || now.Sub(current.startedAt) >= l.window {
		l.windows[key] = &window{startedAt: now, count: 1}
		return true, 0
	}
	if current.count >= l.limit {
		return false, l.window - now.Sub(current.startedAt)
	}
	current.count++
	return true, 0
}

func Middleware(limiter *Limiter, logger *slog.Logger) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, retryAfter := limiter.Allow(session.UserID(r.Context()))
			if !allowed {
				httpx.WriteError(w, r, logger, apperror.NewRetryAfter(apperror.CodeRateLimited, "rate limit exceeded", retryAfter.Milliseconds()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
