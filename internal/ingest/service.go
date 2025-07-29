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
	eventWriter := newBufferedExecutor(s.flushEvents, nil, config.FlushInterval, config.FlushLimit)
	sessionWriter := newBufferedExecutor(s.flushSessions, nil, config.FlushInterval, config.FlushLimit)

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

	return nil
}

func (s *Service) ShutdownWorkers(ctx context.Context) error {
	done := make(chan struct{})	
	go func() {
		defer close(done)
		s.writerWg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (s *Service) IngestEvent(ctx context.Context, event analytics.Event) error {
	if err := s.eventWriter.push(event); err != nil {
		return fmt.Errorf("push event: %w", err)
	}

	return nil
}

func (s *Service) flushEvents(ctx context.Context, events []analytics.Event) error {
	if err := s.storage.BatchInsertEvent(ctx, events); err != nil {
		return fmt.Errorf("batch insert event: %w", err)
	}

	return nil
}

func (s *Service) flushSessions(ctx context.Context, sessions []analytics.Session) error {
	return nil
}
