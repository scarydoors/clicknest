package ingest

import (
	"sync"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
)

type sessionStore struct {
	cache map[string]sessionEntry
	mu sync.Mutex
	ttl time.Duration
}

type sessionEntry struct {
	session analytics.Session
	expiry time.Time
}

func newSessionStore(ttl time.Duration) *sessionStore {
	return &sessionStore{
		cache: make(map[string]sessionEntry),
		ttl: ttl,
	}
}
