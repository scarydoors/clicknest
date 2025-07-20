package clickhouse_test

import (
	"context"
	"testing"
	"time"

	"github.com/scarydoors/clicknest/internal/event"
)

func BenchmarkInsertEvent(b *testing.B) {
	ctx := context.Background()
	evt := event.Event{
		Timestamp: time.Now(),
		Domain: "what.com",
		Kind: "pageview",
		Pathname: "https://what.com/yeah",
	}

	for b.Loop() {
		for range 10000 {
			if err := clickhouseDB.InsertEvent(ctx, evt); err != nil {
				b.Errorf("%s", err)
			}
		}
	}
}

func BenchmarkBatchInsertEvent(b *testing.B) {
	ctx := context.Background()
	evt := event.Event{
		Timestamp: time.Now(),
		Domain: "what.com",
		Kind: "pageview",
		Pathname: "https://what.com/yeah",
	}

	evts := make([]event.Event, 0, 100000)
	for range 100000 {
		evts = append(evts, evt)
	}
	for b.Loop() {
		if err := clickhouseDB.BatchInsertEvent(ctx, evts); err != nil {
			b.Errorf("%s", err)
		}
	}
}
