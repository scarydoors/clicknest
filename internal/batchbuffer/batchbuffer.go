package batchbuffer

import (
	"context"
	"time"

	"github.com/scarydoors/clicknest/internal/workerutil"
	"golang.org/x/sync/singleflight"
)

var (
	_ workerutil.Runner  = (*BatchBuffer[any])(nil)
	_ workerutil.Cleaner = (*BatchBuffer[any])(nil)
)

type FlushConfig struct {
	Interval time.Duration
	Timeout  time.Duration
	Limit    int
}

type BatchBuffer[T any] struct {
	storage       Storage[T]
	errorCallback func(context.Context, error)
	config        FlushConfig

	flushSf singleflight.Group
	itemCh  chan T
	ticker  *time.Ticker
}

type Storage[T any] interface {
	BatchInsert(context.Context, []T) error
}

func NewBatchBuffer[T any](
	storage Storage[T],
	errorCallback func(context.Context, error),
	config FlushConfig,
) *BatchBuffer[T] {
	return &BatchBuffer[T]{
		storage:       storage,
		config:        config,
		errorCallback: errorCallback,

		itemCh: make(chan T, config.Limit),
		ticker: time.NewTicker(config.Interval),
	}
}

func (b *BatchBuffer[T]) Run(ctx context.Context) error {
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

func (b *BatchBuffer[T]) Push(ctx context.Context, item T) error {
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

func (b *BatchBuffer[T]) flush(ctx context.Context) {
	_, _, _ = b.flushSf.Do("flush", func() (any, error) {
		flushContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), b.config.Timeout)
		defer cancel()
		err := b.doFlush(flushContext)
		if err != nil && b.errorCallback != nil {
			b.errorCallback(flushContext, err)
		}
		return nil, nil
	})
}

func (b *BatchBuffer[T]) FinalFlush(ctx context.Context) error {
	_, err, _ := b.flushSf.Do("flush", func() (any, error) {
		err := b.doFlush(ctx)
		b.ticker.Stop()
		return nil, err
	})
	return err
}

func (b *BatchBuffer[T]) Cleanup(ctx context.Context) error {
	return b.FinalFlush(ctx)
}

func (b *BatchBuffer[T]) doFlush(ctx context.Context) error {
	b.ticker.Stop()
	defer b.ticker.Reset(b.config.Interval)

	if len(b.itemCh) == 0 {
		return nil
	}

	count := len(b.itemCh)
	buf := make([]T, 0, count)
	for range count {
		buf = append(buf, <-b.itemCh)
	}

	if err := b.storage.BatchInsert(ctx, buf); err != nil {
		return err
	}

	return nil
}
