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

const DefaultSessionTTL = 30 * time.Second
const DefaultSessionCheckInterval = 1 * time.Minute

type Store struct {
	mutexMap 	  sync.Map
	cache         *cache.Cache[analytics.UserID, State]
	sessionWriter *batchbuffer.BatchBuffer[analytics.Session]

	workerCancel context.CancelFunc
	workerWg     sync.WaitGroup

	logger *slog.Logger
}

type State struct {
	SessionID analytics.SessionID
	Start     time.Time
	End       time.Time
	EventCount uint32
}

func (s State) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("session_id", s.SessionID),
		slog.Time("start", s.Start),
		slog.Time("end", s.End),
		slog.Any("event_count", s.EventCount),
	)
}

func NewStore(config batchbuffer.FlushConfig, storage batchbuffer.Storage[analytics.Session], logger *slog.Logger) *Store {
	s := &Store{
		logger: logger,
	}

	s.sessionWriter = batchbuffer.NewBatchBuffer(storage, s.handleSessionWriterError, config)
	s.cache = cache.NewCache(DefaultSessionTTL, DefaultSessionCheckInterval, s.onSessionExpire)

	return s
}

func (s *Store) ExtendSession(ctx context.Context, event *analytics.Event) error {
	mu := s.getSessionMutex(event.UserID)
	mu.Lock()
	defer mu.Unlock()

	var oldSession analytics.Session
	entry, found := s.cache.Get(event.UserID)
	if !found {
		oldSession = analytics.NewSession(*event)
	} else {
		var err error
		oldSession, err = composeSession(*event, entry.Value)
		if err != nil {
			return err
		}
	}

	newSession, err := oldSession.EventAdded(*event);
	if err != nil {
		return err
	}

	oldSession.MarkCollapse()
	if err := s.sessionWriter.Push(ctx, oldSession); err != nil {
		return err
	}
	if err := s.sessionWriter.Push(ctx, newSession); err != nil {
		return err
	}

	newState := sessionToState(newSession)
	s.logger.Info("state set", slog.Any("state", newState))
	s.cache.Set(event.UserID, newState)

	event.SessionID = newSession.SessionID

	return nil
}

func (s *Store) onSessionExpire(userID analytics.UserID, _ State) {
	// ensure that there is no active update to a session happening.
	mu := s.getSessionMutex(userID)
	mu.Lock()
	defer mu.Unlock()

	// while an update was happening, we could have created a new session for
	// the given userID, we mustn't delete the mutex because it could cause 2
	// parallel updates to session if there is an s.RecordEvent still waiting
	// for the old mutex to be unlocked and an s.RecordEvent creating a new one
	// since it has been deleted from the map
	if _, found := s.cache.Get(userID); !found {
		s.mutexMap.Delete(userID)
	}
}

func (s *Store) getSessionMutex(userID analytics.UserID) *sync.Mutex {
	mu, ok := s.mutexMap.Load(userID)
	if ok {
		return mu.(*sync.Mutex)
	}

	newMu := &sync.Mutex{}
	actual, loaded := s.mutexMap.LoadOrStore(userID, newMu)
	if loaded {
		return actual.(*sync.Mutex)
	}

	return newMu
}

func composeSession(event analytics.Event, state State) (analytics.Session, error) {
	duration, err := analytics.NewSessionDuration(state.Start, state.End)
	if err != nil {
		return analytics.Session{}, fmt.Errorf("calculate session duration: %w", err)
	}

	return analytics.Session{
		Start: state.Start,
		End: state.End,
		Domain: event.Domain,
		Duration: duration,
		EventCount: state.EventCount,
		SessionID: state.SessionID,
		UserID: event.UserID,
		Sign: 1,	
	}, nil
}

func sessionToState(session analytics.Session) State {
	return State{
		SessionID: session.SessionID,
		Start: session.Start,
		End: session.End,
		EventCount: session.EventCount,
	}
}
