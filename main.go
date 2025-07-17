package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/scarydoors/clicknest/internal/server"
)


func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	srv := server.NewServer(logger)

	httpServer := http.Server{
		Addr: ":6969",
		Handler: srv,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		logger.Info("server listening", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("closing server...")
			} else {
				logger.Error("error listening", slog.Any("error", err))
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error while shutting down server", slog.Any("error", err))
		}
	}()

	wg.Wait()
}
