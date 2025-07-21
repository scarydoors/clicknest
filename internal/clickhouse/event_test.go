package clickhouse_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/scarydoors/clicknest/internal/event"
)

func BenchmarkInsertEvent(b *testing.B) {
	ctx := context.Background()
	if err := clickhouseDB.TruncateTables(ctx, "events"); err != nil {
		b.Fatalf("truncate events: %s", err)
	}

	evt := event.Event{
		Timestamp: time.Now(),
		Domain: "what.com",
		Kind: "pageview",
		Pathname: "https://what.com/yeah",
	}

	for b.Loop() {
		if err := clickhouseDB.InsertEvent(ctx, evt); err != nil {
			b.Errorf("%s", err)
		}
	}
	b.ReportMetric(float64(b.N) / b.Elapsed().Seconds(), "events/sec")
}

func BenchmarkAsyncInsertEvent(b *testing.B) {
	ctx := context.Background()
	if err := clickhouseDB.TruncateTables(ctx, "events"); err != nil {
		b.Fatalf("truncate events: %s", err)
	}

	evt := event.Event{
		Timestamp: time.Now(),
		Domain: "what.com",
		Kind: "pageview",
		Pathname: "https://what.com/yeah",
	}

	for b.Loop() {
		if err := clickhouseDB.AsyncInsertEvent(ctx, evt); err != nil {
			b.Errorf("%s", err)
		}
	}
	b.ReportMetric(float64(b.N) / b.Elapsed().Seconds(), "events/sec")
}

func BenchmarkBatchInsertEvent(b *testing.B) {
	ctx := context.Background()
	evt := event.Event{
		Timestamp: time.Now(),
		Domain: "what.com",
		Kind: "pageview",
		Pathname: "https://what.com/yeah",
	}

	batchSizes := [4]int{1000, 10000, 100000, 1000000}
	for _, size := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize-%d", size), func (b *testing.B) {
			if err := clickhouseDB.TruncateTables(ctx, "events"); err != nil {
				b.Fatalf("truncate events: %s", err)
			}
			evts := make([]event.Event, 0, size)
			for range size {
				evts = append(evts, evt)
			}

			for b.Loop() {
				if err := clickhouseDB.BatchInsertEvent(ctx, evts); err != nil {
					b.Errorf("%s", err)
				}
			}

			b.ReportMetric(float64(b.N * size) / b.Elapsed().Seconds(), "events/sec")
		})
	}
}
