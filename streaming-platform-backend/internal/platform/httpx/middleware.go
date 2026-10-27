package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/apperror"
)

type Middleware func(http.Handler) http.Handler

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for index := len(middlewares) - 1; index >= 0; index-- {
		handler = middlewares[index](handler)
	}
	return handler
}

func WithAccessLog(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)
			logger.InfoContext(r.Context(), "request served",
				"requestId", RequestID(r.Context()), "method", r.Method, "path", r.URL.Path,
				"status", recorder.status, "durationMs", time.Since(startedAt).Milliseconds())
		})
	}
}

func WithRecovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer recoverPanic(w, r, logger)
			next.ServeHTTP(w, r)
		})
	}
}

func recoverPanic(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	recovered := recover()
	if recovered == nil {
		return
	}
	logger.ErrorContext(r.Context(), "panic recovered", "requestId", RequestID(r.Context()), "panic", recovered)
	WriteError(w, r, logger, apperror.New(apperror.CodeInternal, "unexpected error"))
}
