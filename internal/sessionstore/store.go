package sessionstore

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/batchbuffer"
	"github.com/scarydoors/clicknest/internal/cache"
	"github.com/scarydoors/clicknest/internal/workerutil"
)

const DefaultSessionTTL = 30 * time.Minute
const DefaultSessionCheckInterval = 1 * time.Minute

type Store struct {
	cache *cache.Cache[analytics.UserID, Entry]
	sessionWriter *batchbuffer.BatchBuffer[analytics.Session]

	workerCancel context.CancelFunc
	workerWg sync.WaitGroup

	logger *slog.Logger
}

type Entry struct {
	SessionID analytics.SessionID
	Domain string
	Start time.Time
	End time.Time
}

func NewStore(config batchbuffer.FlushConfig, storage batchbuffer.Storage[analytics.Session], logger *slog.Logger) *Store {
	s := &Store{
		cache: cache.NewCache[analytics.UserID, Entry](DefaultSessionTTL, DefaultSessionCheckInterval),
		logger: logger,
	}

	s.sessionWriter = batchbuffer.NewBatchBuffer(storage, s.handleSessionWriterError, config)

	return s
}

func (s *Store) Start() error {
	workers := s.workers()

	s.workerCancel = workerutil.StartWorkers(&s.workerWg, s.logger, workers...)
	
	return nil
}

func (s *Store) workers() []workerutil.Worker {
	return []workerutil.Worker{
		{
			Name: "sessionWriter",
			Runner: s.sessionWriter,
		},
		{
			Name: "cache",
			Runner: s.cache,
		},
	}
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
