package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/scarydoors/clicknest/internal/analytics"
)

type Service struct {
	eventStorage Storage[analytics.Event]
	sessionStorage Storage[analytics.Session]
	logger *slog.Logger

	eventWriter *batchBuffer[analytics.Event]
	sessionWriter *batchBuffer[analytics.Session]
	writerCancel context.CancelFunc
	writerWg sync.WaitGroup
}

type Storage[T any] interface {
	BatchInsert(context.Context, []T) error
}

func NewService(eventStorage Storage[analytics.Event], sessionStorage Storage[analytics.Session], logger *slog.Logger) *Service {
	return &Service{
		eventStorage: eventStorage,
		sessionStorage: sessionStorage,
		logger: logger,
	}	
}

type runner interface {
	run(context.Context) error
}

func (s *Service) StartWorkers(config FlushConfig) error {
	eventWriter := newBatchBuffer(s.eventStorage, s.createWriterErrorHandler("event"), config)
	sessionWriter := newBatchBuffer(s.sessionStorage, s.createWriterErrorHandler("session"), config)

	runners := [2]runner{
		eventWriter,
		sessionWriter,
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.writerCancel = cancel

	for _, runner := range runners {
		s.writerWg.Add(1)
		go func() {
			defer s.writerWg.Done()
			runner.run(ctx)
		}()
	}

	s.eventWriter = eventWriter
	s.sessionWriter = sessionWriter

	return nil
}

func (s *Service) ShutdownWorkers(ctx context.Context) error {
	defer func() {
		s.eventWriter = nil
		s.sessionWriter = nil
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

	s.eventWriter.finalFlush(ctx)
	s.sessionWriter.finalFlush(ctx)

	return nil
}

func (s *Service) IngestEvent(ctx context.Context, event analytics.Event) error {
	if s.eventWriter == nil {
		return fmt.Errorf("event worker not running")
	}
	if s.sessionWriter == nil {
		return fmt.Errorf("session worker not running")
	}


	if err := s.eventWriter.push(ctx, event); err != nil {
		return fmt.Errorf("push event: %w", err)
	}

	return nil
}

func (s *Service) createWriterErrorHandler(name string) func(context.Context, error) {
	return func(ctx context.Context, err error) {
		s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("writerName", name), slog.Any("error", err))
	}
}
