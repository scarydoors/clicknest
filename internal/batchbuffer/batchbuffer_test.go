package batchbuffer_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/scarydoors/clicknest/internal/batchbuffer"
	"github.com/stretchr/testify/assert"
)

type testStorage struct {
	batchInsert func(context.Context, []int) error
}

func (ts testStorage) BatchInsert(ctx context.Context, items []int) error {
	return ts.batchInsert(ctx, items)
}

func TestBatchBuffer_FinalFlushClearsOutAllItems(t *testing.T) {
	var flushedItemCount int
	batchInsert := func(ctx context.Context, items []int) error {
		flushedItemCount += len(items)
		return nil
	}

	storage := testStorage{
		batchInsert: batchInsert,
	}

	const limit = 10
	writer := batchbuffer.NewBatchBuffer(storage, nil, batchbuffer.FlushConfig{
		// never flush using the time
		Interval: 24 * time.Hour,
		Timeout:  24 * time.Hour,
		Limit:    limit,
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = writer.Run(ctx)
	}()

	const pushCount = (limit * 3) + limit - 1

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for pushNo := range pushCount {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = writer.Push(context.Background(), pushNo)
			}()
		}
	}()
	wg.Wait()

	cancel()
	<-done
	_ = writer.FinalFlush(context.Background())
	assert.Equal(t, pushCount, flushedItemCount)
}
