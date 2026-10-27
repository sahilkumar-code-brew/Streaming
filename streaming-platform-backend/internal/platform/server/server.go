package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const readHeaderTimeout = 5 * time.Second

func Run(addr string, handler http.Handler, logger *slog.Logger, shutdownTimeout time.Duration) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	httpServer := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: readHeaderTimeout}
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- httpServer.ListenAndServe() }()
	logger.Info("server started", "addr", addr)
	select {
	case err := <-serveErrors:
		return ignoreClosed(err)
	case <-ctx.Done():
		return drain(httpServer, logger, shutdownTimeout)
	}
}

func drain(httpServer *http.Server, logger *slog.Logger, shutdownTimeout time.Duration) error {
	logger.Info("shutdown started", "timeout", shutdownTimeout.String())
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return httpServer.Shutdown(ctx)
}

func ignoreClosed(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
