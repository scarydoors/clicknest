package sessionstore

import (
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/cache"
)

const SessionTtl = 30 * time.Minute
const SessionCheckInterval = 1 * time.Minute

type Store struct {
	cache *cache.Cache[analytics.UserID, Entry]
}

type Entry struct {
	SessionID analytics.SessionID
	Start time.Time
	End time.Time
}


func NewSessionStore() *Store {
	return &Store{
		cache: cache.NewCache[analytics.UserID, Entry](SessionTtl, SessionCheckInterval),
	}
}

func (s *Store) Start() error {
	return s.cache.Start()
}

func (s *Store) Stop() {
	s.cache.Stop()
}
