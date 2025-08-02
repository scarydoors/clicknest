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

type bufferedExecutor[T any] struct {
	config FlushConfig
	execCallback func(context.Context, []T) error
	errorCallback func(context.Context, error)

	flushGroup singleflight.Group
	itemCh chan T
	ticker *time.Ticker
}

func newBufferedExecutor[T any](
	execCallback func(context.Context, []T) error,
	errorCallback func(context.Context, error),
	config FlushConfig,
) *bufferedExecutor[T] {
	return &bufferedExecutor[T]{
		config: config,
		execCallback: execCallback,
		errorCallback: errorCallback,
		itemCh: make(chan T, config.Limit),
		ticker: time.NewTicker(config.Interval),
	}
}

func (b *bufferedExecutor[T]) run(ctx context.Context) error {
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

func (b *bufferedExecutor[T]) push(ctx context.Context, item T) error {
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

func (b *bufferedExecutor[T]) flush(ctx context.Context) {
	b.flushGroup.Do("flush", func() (any, error) {
		flushContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), b.config.Timeout)
		defer cancel()
		return b.doFlush(flushContext)
	})
}

func (b *bufferedExecutor[T]) finalFlush(ctx context.Context) {
	b.flushGroup.Do("flush", func() (any, error) {
		return b.doFlush(ctx)
	})
}

func (b *bufferedExecutor[T]) doFlush(ctx context.Context) (any, error) {
	b.ticker.Stop()
	defer b.ticker.Reset(b.config.Interval)

	buf := make([]T, 0, len(b.itemCh))
	if len(b.itemCh) == 0 {
		return nil, nil
	}

	for len(b.itemCh) > 0 && len(buf) < b.config.Limit {
		buf = append(buf, <-b.itemCh)	
	}

	if err := b.execCallback(ctx, buf); err != nil {
		b.errorCallback(ctx, err)
	}

	return nil, nil
}
