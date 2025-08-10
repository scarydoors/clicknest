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

func NewStore(config batchbuffer.FlushConfig, storage batchbuffer.Storage[analytics.Session], logger *slog.Logger) *Store {
	s := &Store{
		cache:  cache.NewCache[analytics.UserID, State](DefaultSessionTTL, DefaultSessionCheckInterval),
		logger: logger,
	}

	s.sessionWriter = batchbuffer.NewBatchBuffer(storage, s.handleSessionWriterError, config)

	return s
}

func (s *Store) RecordEvent(ctx context.Context, event *analytics.Event) error {
	var oldSession analytics.Session
	entry, found := s.cache.Get(event.UserID)
	if !found {
		oldSession = analytics.NewSession(*event)
	} else {
		oldSession = composeSession(*event, entry.Value)
	}

	newSession, err := oldSession.EventAdded(*event);
	if err != nil {
		return err
	}

	oldSession.MarkCollapse()
	s.sessionWriter.Push(ctx, oldSession)
	s.sessionWriter.Push(ctx, newSession)

	newState := sessionToState(newSession)
	s.cache.Set(event.UserID, newState)

	event.SessionID = newSession.SessionID

	return nil
}

func composeSession(event analytics.Event, state State) analytics.Session {
	duration, err := analytics.NewSessionDuration(state.Start, state.End)
	if err != nil {
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
	}
}

func sessionToState(session analytics.Session) State {
	return State{
		SessionID: session.SessionID,
		Start: session.Start,
		End: session.End,
		EventCount: session.EventCount,
	}
}
