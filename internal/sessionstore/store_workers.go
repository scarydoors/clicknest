package sessionstore

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/scarydoors/clicknest/internal/workerutil"
)

func (s *Store) workers() []workerutil.Worker {
	return []workerutil.Worker{
		{
			Name:   "sessionWriter",
			Runner: s.sessionWriter,
		},
		{
			Name:   "cache",
			Runner: s.cache,
		},
	}
}

func (s *Store) Start() error {
	workers := s.workers()

	s.workerCancel = workerutil.StartWorkers(&s.workerWg, s.logger, workers...)

	return nil
}

func (s *Store) handleSessionWriterError(ctx context.Context, err error) {
	s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("name", "session"), slog.Any("error", err))
}

func (s *Store) Shutdown(ctx context.Context) error {
	s.workerCancel()

	workers := s.workers()

	if err := workerutil.CleanupWorkers(ctx, &s.workerWg, s.logger, workers...); err != nil {
		return fmt.Errorf("cleanup workers: %w", err)
	}

	return nil
}
