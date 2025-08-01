package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
)

type Service struct {
	storage Storage
	logger *slog.Logger

	writerCancel context.CancelFunc
	eventWriter *bufferedExecutor[analytics.Event]
	sessionWriter *bufferedExecutor[analytics.Session]
	writerWg sync.WaitGroup
}

type Storage interface {
	BatchInsertEvent(context.Context, []analytics.Event) error
}

func NewService(storage Storage, logger *slog.Logger) *Service {
	return &Service{
		storage: storage,
		logger: logger,
	}	
}

type WorkerConfig struct {
	FlushInterval time.Duration
	FlushLimit int
}

type runner interface {
	run(context.Context) error
}

func (s *Service) StartWorkers(config WorkerConfig) error {
	eventWriter := newBufferedExecutor(s.handleEventFlush, s.createWriterErrorHandler("event"), config.FlushInterval, config.FlushLimit)
	sessionWriter := newBufferedExecutor(s.handleSessionFlush, s.createWriterErrorHandler("session"), config.FlushInterval, config.FlushLimit)

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

	s.eventWriter.flush(ctx)
	s.sessionWriter.flush(ctx)

	return nil
}

func (s *Service) IngestEvent(ctx context.Context, event analytics.Event) error {
	if s.eventWriter == nil {
		return fmt.Errorf("event worker not running")
	}

	if err := s.eventWriter.push(event); err != nil {
		return fmt.Errorf("push event: %w", err)
	}

	return nil
}

func (s *Service) handleEventFlush(ctx context.Context, events []analytics.Event) error {
	if err := s.storage.BatchInsertEvent(ctx, events); err != nil {
		return err
	}

	return nil
}

func (s *Service) handleSessionFlush(ctx context.Context, sessions []analytics.Session) error {
	if s.sessionWriter == nil {
		return fmt.Errorf("session worker not running")
	}

	return nil
}

func (s *Service) createWriterErrorHandler(name string) func(context.Context, error) {
	return func(ctx context.Context, err error) {
		s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("writerName", name), slog.Any("error", err))
	}
}
