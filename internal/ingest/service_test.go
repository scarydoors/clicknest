package ingest_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {}

var processedItemsNormally bool
var finalFlushExecuted bool
var batchInsertTimesCalled int

func (s *mockStorage) BatchInsertEvent(ctx context.Context, _ []analytics.Event) error {
	batchInsertTimesCalled++
	_, ok := ctx.Deadline()
	if ok {
		finalFlushExecuted = true
	} else {
		processedItemsNormally = true
	}

	return nil
}

func TestServiceWorkersLifecycle(t *testing.T) {
	service := ingest.NewService(&mockStorage{}, slog.Default())
	service.StartWorkers(ingest.WorkerConfig{
		FlushInterval: 60 * time.Second,
		FlushLimit: 2,
	})

	for range 3 {
		service.IngestEvent(context.Background(), analytics.Event{})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	service.ShutdownWorkers(ctx)

	assert.Equal(t, 2, batchInsertTimesCalled)
	assert.True(t, processedItemsNormally)
	assert.True(t, finalFlushExecuted)
}
