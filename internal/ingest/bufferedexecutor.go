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
	readyCh chan struct{}

	ctx context.Context
}

func newBufferedExecutor[T any](
	execCallback func(context.Context, []T) error,
	errorCallback func(error),
	flushInterval time.Duration,
	flushLimit int,
) *bufferedExecutor[T] {
	return &bufferedExecutor[T]{
		flushInterval: flushInterval,
		flushLimit: flushLimit,
		execCallback: execCallback,
		errorCallback: errorCallback,

		buf: make([]T, 0, flushLimit),
		itemCh: make(chan T),
		readyCh: make(chan struct{}),
	}
}

func (b *bufferedExecutor[T]) run(ctx context.Context) error {
	tickCh := time.Tick(b.flushInterval)
	defer close(b.itemCh)

	b.ctx = ctx
	close(b.readyCh)

	for {
		select {
		case <-b.ctx.Done():
			return b.ctx.Err()

		case <-tickCh:
			b.flush(context.Background())

		case item := <-b.itemCh:
			b.buf = append(b.buf, item)
			if len(b.buf) >= b.flushLimit {
				b.flush(context.Background())
			}
		}
	}
}

func (b *bufferedExecutor[T]) push(item T) error {
	<-b.readyCh

	select {
	case <-b.ctx.Done():
		return b.ctx.Err()
	case b.itemCh <- item:
		return nil
	}
}

func (b *bufferedExecutor[T]) flush(ctx context.Context) {
	if len(b.buf) == 0 {
		return
	}

	if err := b.execCallback(ctx, b.buf); err != nil {
		b.errorCallback(err)
	}

	b.buf = b.buf[:0]
}
