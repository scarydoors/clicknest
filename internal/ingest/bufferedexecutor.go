package ingest

import (
	"context"
	"time"
)

type bufferedExecutor[T any] struct {
	flushInterval time.Duration
	flushLimit int
	execCallback func(context.Context, []T) error
	errorCallback func(error)

	buf []T
	itemCh chan T

	ctx context.Context
}

func newBufferedExecutor[T any](
	execCallback func(context.Context, []T) error,
	errorCallback func(error),
	flushInterval time.Duration,
	flushLimit int,
) *bufferedExecutor[T] {
	return &bufferedExecutor[T]{
		buf: make([]T, 0, flushLimit),
		flushInterval: flushInterval,
		flushLimit: flushLimit,
		execCallback: execCallback,
		errorCallback: errorCallback,
	}
}

func (b *bufferedExecutor[T]) run(ctx context.Context) error {
	tickCh := time.Tick(b.flushInterval)

	b.itemCh = make(chan T)
	defer close(b.itemCh)

	b.ctx = ctx

	for {
		select {
		case <-b.ctx.Done():
			if len(b.buf) > 0 {
				b.flush()
			}
			return b.ctx.Err()

		case <-tickCh:
			if len(b.buf) > 0 {
				b.flush()
			}

		case item := <-b.itemCh:
			b.buf = append(b.buf, item)
			if len(b.buf) >= b.flushLimit {
				b.flush()
			}
		}
	}
}

func (b *bufferedExecutor[T]) push(item T) error {
	select {
	case <-b.ctx.Done():
		return b.ctx.Err()
	case b.itemCh <- item:
		return nil
	}
}

func (b *bufferedExecutor[T]) flush() {
	ctx := context.Background()

	if err := b.execCallback(ctx, b.buf); err != nil {
		b.errorCallback(err)
	}

	b.buf = b.buf[:0]
}
