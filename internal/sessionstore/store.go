package sessionstore

import (
	"context"
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
	cache         *cache.Cache[analytics.UserID, Entry]
	sessionWriter *batchbuffer.BatchBuffer[analytics.Session]

	workerCancel context.CancelFunc
	workerWg     sync.WaitGroup

	logger *slog.Logger
}

type Entry struct {
	SessionID analytics.SessionID
	Domain    string
	Start     time.Time
	End       time.Time
}

func NewStore(config batchbuffer.FlushConfig, storage batchbuffer.Storage[analytics.Session], logger *slog.Logger) *Store {
	s := &Store{
		cache:  cache.NewCache[analytics.UserID, Entry](DefaultSessionTTL, DefaultSessionCheckInterval),
		logger: logger,
	}

	s.sessionWriter = batchbuffer.NewBatchBuffer(storage, s.handleSessionWriterError, config)

	return s
}
