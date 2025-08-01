package ingest_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/ingest"
	//"github.com/stretchr/testify/assert"
)

type mockStorage struct {}

func (s *mockStorage) BatchInsertEvent(ctx context.Context, events []analytics.Event) error {
	log.Default().Print("processing", events, "\n")
	if _, ok := ctx.Deadline(); ok {
		log.Default().Print("cleaning", "\n")
	} else if len(events) == 1 {
		log.Default().Print("handleing 1 event case")
		select {
		case <-time.After(60 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	} else {
		log.Default().Print("normal run")
	}
	log.Default().Print("____ ")
	return nil;
}

func TestServiceWorkersLifecycle(t *testing.T) {
	service := ingest.NewService(&mockStorage{}, slog.Default())
	service.StartWorkers(ingest.WorkerConfig{
		FlushInterval: 2 * time.Second,
		FlushLimit: 2,
	})

	for n := range 3 {
		service.IngestEvent(context.Background(), analytics.Event{Domain: fmt.Sprintf("%d", n)})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	time.Sleep(5 * time.Second)
	service.ShutdownWorkers(ctx)
}
