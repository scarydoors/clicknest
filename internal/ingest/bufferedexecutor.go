package ingest

import (
	"context"
	"time"
)

type bufferedExecutor[T any] struct {
	buf []T
	flushInterval time.Duration
	flushLimit int
	execCallback func([]T) error
	errorCallback func(error)
	itemCh chan T
}

func newBufferedExecutor[T any](
	execCallback func([]T) error,
	errorCallback func(error),
	flushInterval time.Duration,
	flushLimit int,
) bufferedExecutor[T] {
	return bufferedExecutor[T]{
		buf: make([]T, 0, flushLimit),
		flushInterval: flushInterval,
		flushLimit: flushLimit,
		execCallback: execCallback,
		errorCallback: errorCallback,
	}
}

func (b *bufferedExecutor[T]) run(ctx context.Context) error {
	b.itemCh = make(chan T)
	tickerCh := time.NewTicker(b.flushInterval)

	for {
		select {
		case <-ctx.Done():

		case <-tickerCh.C:
			if (len(b.buf) >= b.flushLimit) {
			}
		}
	}
}

func (b *bufferedExecutor[T]) push(item T) {
	b.itemCh <- item
}
