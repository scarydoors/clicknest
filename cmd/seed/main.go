package main

import (
	"context"
	"math/rand/v2"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
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

	eventRepo := clickhouse.NewEventRepository(clickhouseDB, logger)
	_ = clickhouse.NewSessionRepository(clickhouseDB, logger)

	if err := clickhouseDB.Exec(ctx, "TRUNCATE ALL TABLES FROM default"); err != nil {
		slog.Error("failed clickhouse truncate", slog.Any("error", err))
		os.Exit(1)
	}

	startPage := "/home"
	endPage := "/checkout"
	pages := []string{
		"/product/1",
		"/product/2",
		"/product/3",
		"/product/4",
		"/product/5",
		"/product/6",
		"/product/7",
	}

	domain := "stupidwebsite.com"
	fullPathFn := func (pathname string) string { return fmt.Sprintf("http://%s%s", domain, pathname) }

	initialTimestamp := time.Now().Add(time.Duration(-24 * time.Hour))
	events := []analytics.Event{}
	// 3 users
	for range(3) {
		userID := analytics.UserID(rand.Uint64())
		timestamp := initialTimestamp
		startEv, err := analytics.NewEvent(timestamp, domain, analytics.EventKindPageview, fullPathFn(startPage))
		startEv.UserID = userID
		if err != nil {
			panic(err)
		}
		events = append(events, startEv)

		for range(rand.Int32N(500)) {
			timestamp := timestamp.Add(time.Duration(rand.Uint64N(30) * uint64(time.Minute)))
			ev, err := analytics.NewEvent(timestamp, domain, analytics.EventKindPageview, fullPathFn(pages[rand.IntN(len(pages))]))
			ev.UserID = userID
			if err != nil {
				panic(err)
			}
			events = append(events, ev)
		}

		endEv, err := analytics.NewEvent(timestamp, domain, analytics.EventKindPageview, fullPathFn(endPage))
		if err != nil {
			panic(err)
		}
		endEv.UserID = userID
		events = append(events, endEv)
	}

	if err := eventRepo.BatchInsert(ctx, events); err != nil {
		logger.Error("batch insert events", slog.Any("error", err))
		os.Exit(1)
	}
}
