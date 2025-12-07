package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/scarydoors/clicknest/internal/clickhouse"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	clickhouseDB, err := clickhouse.NewClickhouseConn(ctx, clickhouse.ClickhouseDBConfig{
		Host:     "localhost",
		Port:     "9000",
		Database: "default",
		Username: "default",
		Password: "",
	})

	if err != nil {
		slog.Error("failed clickhouse init", slog.Any("error", err))
		os.Exit(1)
	}

	_ = clickhouse.NewEventRepository(clickhouseDB, logger)
	_ = clickhouse.NewSessionRepository(clickhouseDB, logger)

	if err := clickhouseDB.Exec(ctx, "TRUNCATE ALL TABLES FROM default"); err != nil {
		slog.Error("failed clickhouse truncate", slog.Any("error", err))
		os.Exit(1)
	}
}
