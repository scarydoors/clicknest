package workerutil

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Runner interface {
	Run(context.Context) error
}

type Cleaner interface {
	Cleanup(context.Context) error
}

type Worker struct {
	Name   string
	Runner Runner
}

func StartWorkers(wg *sync.WaitGroup, logger *slog.Logger, workers ...Worker) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	for _, worker := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// only returns context error, error is ignored
			logger.Info("starting worker", slog.String("name", worker.Name))
			if err := worker.Runner.Run(ctx); err == context.Canceled {
				logger.Info("worker stopped running", slog.String("name", worker.Name))
			} else if err != nil {
				logger.Error("worker exited with error", slog.String("name", worker.Name), slog.Any("error", err))
			}
		}()
	}

	return cancel
}

func CleanupWorkers(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger, workers ...Worker) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("unable to shutdown gracefully: %w", ctx.Err())
	case <-done:
	}

	var cleanupWg sync.WaitGroup
	for _, worker := range workers {
		c, ok := worker.Runner.(Cleaner)
		if !ok {
			continue
		}

		cleanupWg.Add(1)
		go func() {
			defer cleanupWg.Done()
			if err := c.Cleanup(ctx); err != nil {
				logger.Error("worker cleanup encountered an error", slog.String("name", worker.Name), slog.Any("error", err))
			} else {
				logger.Info("worker cleanup finished", slog.String("name", worker.Name))
			}
		}()
	}

	return nil
}
