package sessionstore

import (
	"time"

	"github.com/scarydoors/clicknest/internal/cache"
)

const SessionTtl = 30 * time.Minute
const SessionCheckInterval = 1 * time.Minute

type Store struct {
	cache *cache.Cache[uint64, Entry]
}

type Entry struct {
	sessionID uint64
	start time.Time
	end time.Time
}


func NewSessionStore() *Store {
	return &Store{
		cache: cache.NewCache[uint64, Entry](SessionTtl, SessionCheckInterval),
	}
}

func (s *Store) Start() error {
	return s.cache.Start()
}

func (s *Store) Stop() {
	s.cache.Stop()
}
