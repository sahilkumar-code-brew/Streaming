package httpx

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-Id"

type requestIDKey struct{}

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := resolveRequestID(r.Header.Get(HeaderRequestID))
		w.Header().Set(HeaderRequestID, requestID)
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}

func resolveRequestID(incoming string) string {
	if _, err := uuid.Parse(incoming); err == nil {
		return incoming
	}
	return uuid.NewString()
}
