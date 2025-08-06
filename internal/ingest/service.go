package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/batchbuffer"
	"github.com/scarydoors/clicknest/internal/sessionstore"
)

type Service struct {
	logger         *slog.Logger
	sessionStore *sessionstore.Store
	config batchbuffer.FlushConfig

	eventStorage   batchbuffer.Storage[analytics.Event]
	eventWriter   *batchbuffer.BatchBuffer[analytics.Event]
	writerCancel  context.CancelFunc
	writerWg      sync.WaitGroup
}

func NewService(config batchbuffer.FlushConfig, eventStorage batchbuffer.Storage[analytics.Event], sessionStore *sessionstore.Store, logger *slog.Logger) *Service {
	return &Service{
		config: config,
		eventStorage:   eventStorage,
		sessionStore: sessionStore,
		logger:         logger,
	}
}

type writer interface {
	Run(context.Context) error
	FinalFlush(context.Context) error
}

func (s *Service) Start() error {
	eventWriter := batchbuffer.NewBatchBuffer(s.eventStorage, s.handleEventWriterError, s.config)

	ctx, cancel := context.WithCancel(context.Background())
	s.writerCancel = cancel

	s.writerWg.Add(1)
	go func() {
		defer s.writerWg.Done()
		// only returns context error, error is ignored
		if err := eventWriter.Run(ctx); err == context.Canceled {
			s.logger.Info("writer stopped running", slog.String("name", "event"))
		} else if err != nil {
			s.logger.Error("writer exited with error", slog.String("name", "event"), slog.Any("error", err))
		}
	}()

	s.eventWriter = eventWriter

	return nil
}

func (s *Service) ShutdownWorkers(ctx context.Context) error {
	defer func() {
		s.eventWriter = nil
	}()

	s.writerCancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.writerWg.Wait()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("unable to perform final flush: %w", ctx.Err())
	case <-done:
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.eventWriter.FinalFlush(ctx); err != nil {
			s.logger.Error("worker final flush encountered an error", slog.String("name", "event"), slog.Any("error", err))
		} else {
			s.logger.Info("worker final flush gracefully completed", slog.String("name", "event"))
		}
	}()
	wg.Wait()

	return nil
}

func (s *Service) IngestEvent(ctx context.Context, event analytics.Event) error {
	if s.eventWriter == nil {
		return fmt.Errorf("event worker not running")
	}

	if err := s.eventWriter.Push(ctx, event); err != nil {
		return fmt.Errorf("push event: %w", err)
	}

	return nil
}

func (s *Service) handleEventWriterError(ctx context.Context, err error) {
	s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("name", "event"), slog.Any("error", err))
}
