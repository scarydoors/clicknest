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
