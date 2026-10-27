package health

import (
	"context"
	"net/http"
	"time"
)

const readinessTimeout = 2 * time.Second

type Check func(ctx context.Context) error

func Register(mux *http.ServeMux, checks ...Check) {
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(readinessStatus(r.Context(), checks))
	})
}

func readinessStatus(parent context.Context, checks []Check) int {
	ctx, cancel := context.WithTimeout(parent, readinessTimeout)
	defer cancel()
	for _, check := range checks {
		if check(ctx) != nil {
			return http.StatusServiceUnavailable
		}
	}
	return http.StatusOK
}
