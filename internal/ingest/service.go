package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/batchbuffer"
	"github.com/scarydoors/clicknest/internal/sessionstore"
	"github.com/scarydoors/clicknest/internal/workerutil"
)

type Service struct {
	logger       *slog.Logger
	sessionStore *sessionstore.Store

	eventWriter  *batchbuffer.BatchBuffer[analytics.Event]
	workerCancel context.CancelFunc
	workerWg     sync.WaitGroup
}

func NewService(config batchbuffer.FlushConfig, eventStorage batchbuffer.Storage[analytics.Event], sessionStore *sessionstore.Store, logger *slog.Logger) *Service {
	s := &Service{
		logger:       logger,
		sessionStore: sessionStore,
	}

	s.eventWriter = batchbuffer.NewBatchBuffer(eventStorage, s.handleEventWriterError, config)

	return s
}

func (s *Service) Start() error {
	worker := workerutil.Worker{
		Name:   "eventWriter",
		Runner: s.eventWriter,
	}

	s.workerCancel = workerutil.StartWorkers(&s.workerWg, s.logger, worker)

	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.workerCancel()

	worker := workerutil.Worker{
		Name:   "eventWriter",
		Runner: s.eventWriter,
	}

	if err := workerutil.CleanupWorkers(ctx, &s.workerWg, s.logger, worker); err != nil {
		return fmt.Errorf("cleanup workers: %w", err)
	}

	return nil
}

func (s *Service) IngestEvent(ctx context.Context, event analytics.Event) error {
	if s.eventWriter == nil {
		return fmt.Errorf("event worker not running")
	}

	if err := s.sessionStore.ExtendSession(ctx, &event); err != nil {
		return fmt.Errorf("extend session: %w", err)
	}

	if err := s.eventWriter.Push(ctx, event); err != nil {
		return fmt.Errorf("push event: %w", err)
	}

	return nil
}

func (s *Service) handleEventWriterError(ctx context.Context, err error) {
	s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("name", "event"), slog.Any("error", err))
}
