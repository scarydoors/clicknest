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
)

const DefaultSessionTTL = 30 * time.Minute
const DefaultSessionCheckInterval = 1 * time.Minute

type Store struct {
	config batchbuffer.FlushConfig

	cache *cache.Cache[analytics.UserID, Entry]

	sessionStorage batchbuffer.Storage[analytics.Session]
	sessionWriter *batchbuffer.BatchBuffer[analytics.Session]
	writerCancel context.CancelFunc
	writerWg sync.WaitGroup

	logger *slog.Logger
}

type Entry struct {
	SessionID analytics.SessionID
	Domain string
	Start time.Time
	End time.Time
}

func NewStore(config batchbuffer.FlushConfig, sessionStorage batchbuffer.Storage[analytics.Session], logger *slog.Logger) *Store {
	return &Store{
		config: config,
		cache: cache.NewCache[analytics.UserID, Entry](DefaultSessionTTL, DefaultSessionCheckInterval),
		sessionStorage: sessionStorage,
		logger: logger,
	}
}

func (s *Store) Start() error {
	err := s.cache.Start()
	if err != nil {
		return fmt.Errorf("cache start: %w", err)
	}

	sessionWriter := batchbuffer.NewBatchBuffer(s.sessionStorage, s.handleSessionWriterError, s.config)

	ctx, cancel := context.WithCancel(context.Background())
	s.writerCancel = cancel

	s.writerWg.Add(1)
	go func() {
		defer s.writerWg.Done()
		// only returns context error, error is ignored
		if err := sessionWriter.Run(ctx); err == context.Canceled {
			s.logger.Info("writer stopped running", slog.String("name", "event"))
		} else if err != nil {
			s.logger.Error("writer exited with error", slog.String("name", "event"), slog.Any("error", err))
		}
	}()

	s.sessionWriter = sessionWriter

	return nil
}

func (s *Store) handleSessionWriterError(ctx context.Context, err error) {
	s.logger.ErrorContext(ctx, "failed to flush writer", slog.String("name", "session"), slog.Any("error", err))
}

func (s *Store) Stop() {
	s.cache.Stop()
}
