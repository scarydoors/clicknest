package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/scarydoors/clicknest/internal/clickhouse"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/server"
)


func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	clickhouseDB, err := clickhouse.NewClickhouseConn(ctx, clickhouse.ClickhouseDBConfig{
		Host: "localhost",
		Port: "9000",
		Database: "default",
		Username: "default",
		Password: "",
	})

	if err != nil {
		log.Fatalf("failed clickhouse init: %s", err)
	}

	defer clickhouseDB.Close()

	eventRepo := clickhouse.NewEventRepository(clickhouseDB)

	ingestService := ingest.NewService(eventRepo, logger)
	srv := server.NewServer(logger, ingestService)

	httpServer := http.Server{
		Addr: ":6969",
		Handler: srv,
	}

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
