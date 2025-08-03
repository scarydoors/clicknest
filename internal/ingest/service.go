package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/scarydoors/clicknest/internal/analytics"
)

type Service struct {
	eventStorage   Storage[analytics.Event]
	sessionStorage Storage[analytics.Session]
	logger         *slog.Logger

	eventWriter   *batchBuffer[analytics.Event]
	sessionWriter *batchBuffer[analytics.Session]
	writerCancel  context.CancelFunc
	writerWg      sync.WaitGroup
}

type Storage[T any] interface {
	BatchInsert(context.Context, []T) error
}

func NewService(eventStorage Storage[analytics.Event], sessionStorage Storage[analytics.Session], logger *slog.Logger) *Service {
	return &Service{
		eventStorage:   eventStorage,
		sessionStorage: sessionStorage,
		logger:         logger,
	}
}

const eventWriterName = "event"
const sessionWriterName = "session"

type writer interface {
	run(context.Context) error
	finalFlush(context.Context) error
}

type worker struct {
	name   string
	writer writer
}

func (s *Service) StartWorkers(config FlushConfig) error {
	eventWriter := newBatchBuffer(s.eventStorage, s.createWriterErrorHandler(eventWriterName), config)
	sessionWriter := newBatchBuffer(s.sessionStorage, s.createWriterErrorHandler(sessionWriterName), config)

	workers := []worker{
		{eventWriterName, eventWriter},
		{sessionWriterName, sessionWriter},
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.writerCancel = cancel

	for _, worker := range workers {
		s.writerWg.Add(1)
		go func() {
			defer s.writerWg.Done()
			// only returns context error, error is ignored
			if err := worker.writer.run(ctx); err == context.Canceled {
				s.logger.Info("writer stopped running", slog.String("name", worker.name))
			} else if err != nil {
				s.logger.Error("writer exited with error", slog.String("name", worker.name), slog.Any("error", err))
			}
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

	workers := []worker{
		{eventWriterName, s.eventWriter},
		{sessionWriterName, s.sessionWriter},
	}

	var wg sync.WaitGroup
	for _, worker := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := worker.writer.finalFlush(ctx); err != nil {
				s.logger.Error("worker final flush encountered an error", slog.String("name", worker.name), slog.Any("error", err))
			} else {
				s.logger.Info("worker final flush gracefully completed", slog.String("name", worker.name))
			}
		}()
	}
	wg.Wait()

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
		s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("name", name), slog.Any("error", err))
	}
}
