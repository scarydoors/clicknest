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
	buf []T
	itemCh chan T
	timer <-chan time.Time
	readyCh chan struct{}

	ctx context.Context
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

		buf: make([]T, 0, config.Limit),
		itemCh: make(chan T, config.Limit),
		readyCh: make(chan struct{}),
	}
}

func (b *bufferedExecutor[T]) run(ctx context.Context) error {
	defer close(b.itemCh)

	b.resetTimer()

	b.ctx = ctx
	close(b.readyCh)

	for {
		select {
		case <-b.ctx.Done():
			return b.ctx.Err()

		case <-b.timer:
			flushCtx, cancel := b.flushContext()
			defer cancel()
			b.flush(flushCtx)
		}
	}
}

func (b *bufferedExecutor[T]) resetTimer() {
	b.timer = time.After(b.config.Interval)
}

func (b *bufferedExecutor[T]) push(item T) error {
	<-b.readyCh

	for {
		select {
		case <-b.ctx.Done():
			return b.ctx.Err()
		case b.itemCh <- item:
			return nil
		default:
			ctx, cancel := b.flushContext()
			defer cancel()
			b.flush(ctx)
		}
	}
}

func (b *bufferedExecutor[T]) flushContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(b.ctx), b.config.Timeout)
}

func (b *bufferedExecutor[T]) flush(ctx context.Context) {
	b.flushGroup.Do("flush", func() (any, error) {
		defer b.resetTimer()

		if len(b.itemCh) == 0 {
			return nil, nil
		}

		for len(b.itemCh) > 0 && len(b.buf) < b.config.Limit {
			b.buf = append(b.buf, <-b.itemCh)	
		}

		if err := b.execCallback(ctx, b.buf); err != nil {
			b.errorCallback(ctx, err)
		}

		b.buf = b.buf[:0]
		return nil, nil
	})
}
