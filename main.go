package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/scarydoors/clicknest/internal/batchbuffer"
	"github.com/scarydoors/clicknest/internal/clickhouse"
	"github.com/scarydoors/clicknest/internal/errorutil"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/postgres"
	"github.com/scarydoors/clicknest/internal/server"
	"github.com/scarydoors/clicknest/internal/sessionstore"
	"github.com/scarydoors/clicknest/internal/stats"
	"github.com/scarydoors/clicknest/internal/users"
	"github.com/scarydoors/clicknest/internal/validatorutil"
	"github.com/scarydoors/clicknest/internal/workerutil"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := clickhouse.ParseDSN(os.Getenv("CLICKHOUSE_DB_DSN"))
	if err != nil {
		log.Fatalf("failed clickhouse init: %s", err)
	}

	clickhouseDB, err := clickhouse.NewClickhouseConn(ctx, config)

	if err != nil {
		log.Fatalf("failed clickhouse init: %s", err)
	}

	defer errorutil.DeferIgnoreErr(clickhouseDB.Close)

	postgresDB, err := postgres.NewPostgresConn(ctx, os.Getenv("POSTGRES_DB_DSN"))

	flushConfig := batchbuffer.FlushConfig{
		Interval: 4 * time.Second,
		Limit:    100000,
		Timeout:  10 * time.Second,
	}

	validate := validator.New()
	validatorutil.SetupCustomValidations(validate, logger)

	eventRepo := clickhouse.NewEventRepository(clickhouseDB, logger)
	sessionRepo := clickhouse.NewSessionRepository(clickhouseDB, logger)
	statsRepo := clickhouse.NewStatsRepository(clickhouseDB, logger)
	userRepo := postgres.NewUserRepository(postgresDB, logger)

	sessionStore := sessionstore.NewStore(flushConfig, sessionRepo, logger)
	ingestService := ingest.NewService(flushConfig, eventRepo, sessionStore, logger)
	statsService := stats.NewService(statsRepo, logger, validate)
	// TODO: actually use this
	_ = users.NewService(userRepo, logger)

	if err := ingestService.Start(); err != nil {
		log.Fatalf("unable to start ingest service workers: %s", err)
	}
	if err := sessionStore.Start(); err != nil {
		log.Fatalf("unable to start session store workers: %s", err)
	}

	srv := server.NewServer(logger, validate, ingestService, statsService)

	httpServer := http.Server{
		Addr:    ":6969",
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

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error while shutting down server", slog.Any("error", err))
		}

		services := []workerutil.Service{
			{
				Name:       "ingest",
				Shutdowner: ingestService,
			},
			{
				Name:       "sessionstore",
				Shutdowner: sessionStore,
			},
		}
		if err := workerutil.ShutdownServices(shutdownCtx, services...); err != nil {
			for _, err := range errorutil.IntoSlice(err) {
				var serr *workerutil.ShutdownError
				if errors.As(err, &serr) {
					logger.Error("error while shutting down service", slog.String("service", serr.Name), slog.Any("error", serr.Error()))
				}
			}
		}
	}()
	<-done
}
