package ingest

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

type FlushConfig struct {
	Interval time.Duration
	Timeout time.Duration
	Limit int
}

type batchBuffer[T any] struct {
	storage Storage[T]
	errorCallback func(context.Context, error)
	config FlushConfig

	flushGroup singleflight.Group
	itemCh chan T
	ticker *time.Ticker
}

func newBatchBuffer[T any](
	storage Storage[T],
	errorCallback func(context.Context, error),
	config FlushConfig,
) *batchBuffer[T] {
	return &batchBuffer[T]{
		storage: storage,
		config: config,
		errorCallback: errorCallback,

		itemCh: make(chan T, config.Limit),
		ticker: time.NewTicker(config.Interval),
	}
}

func (b *batchBuffer[T]) run(ctx context.Context) error {
	defer close(b.itemCh)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-b.ticker.C:
			b.flush(ctx)
		}
	}
}

func (b *batchBuffer[T]) push(ctx context.Context, item T) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b.itemCh <- item:
			return nil
		default:
			b.flush(ctx)
		}
	}
}

func (b *batchBuffer[T]) flush(ctx context.Context) {
	b.flushGroup.Do("flush", func() (any, error) {
		flushContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), b.config.Timeout)
		defer cancel()
		return b.doFlush(flushContext)
	})
}

func (b *batchBuffer[T]) finalFlush(ctx context.Context) {
	b.flushGroup.Do("flush", func() (any, error) {
		return b.doFlush(ctx)
	})
}

func (b *batchBuffer[T]) doFlush(ctx context.Context) (any, error) {
	b.ticker.Stop()
	defer b.ticker.Reset(b.config.Interval)

	if len(b.itemCh) == 0 {
		return nil, nil
	}

	count := min(len(b.itemCh), b.config.Limit)
	buf := make([]T, 0, count)
	for range count {
		buf = append(buf, <-b.itemCh)	
	}

	if err := b.storage.BatchInsert(ctx, buf); err != nil {
		b.errorCallback(ctx, err)
	}

	return nil, nil
}
