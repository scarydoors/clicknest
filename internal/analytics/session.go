package analytics

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

type SessionID uint64

func NewSessionID() SessionID {
	var id uint64
	// 0 is an invalid id here
	for id == 0 {
		id = rand.Uint64()
	}

	return SessionID(id)
}

type Session struct {
	Start      time.Time
	End        time.Time
	Domain     string
	Duration   SessionDuration
	EventCount uint32
	SessionID  SessionID
	UserID     UserID
	Sign       int8
}

var ErrNegativeDuration = errors.New("session duration cannot be negative")
var ErrDurationOverflowed = errors.New("session duration has overflowed")

type SessionDuration uint32
func NewSessionDuration(start, end time.Time) (SessionDuration, error) {
	duration := end.Sub(start)
	if duration < 0 {
		return 0, ErrNegativeDuration
	}

	secs := duration.Seconds()
	if secs > math.MaxUint32 {
		return 0, ErrDurationOverflowed
	}
	return SessionDuration(secs), nil
}

func (s SessionDuration) Uint32() uint32 {
	return uint32(s)
}

func (s SessionDuration) Duration() time.Duration {
	return time.Duration(s.Uint32()) * time.Second
}
	
func NewSession(event Event) Session {
	return Session{
		Start:      event.Timestamp,
		End:        event.Timestamp,
		Domain:     event.Domain,
		Duration:   0,
		EventCount: 1,
		SessionID:  NewSessionID(),
		UserID:     event.UserID,
		Sign:       1,
	}
}

func (s Session) EventAdded(event Event) (Session, error) {
	if event.Timestamp.After(s.End) {
		duration, err := NewSessionDuration(s.Start, event.Timestamp)
		if err != nil {
			return Session{}, fmt.Errorf("update session: %w", err)
		}

		s.End = event.Timestamp
		s.Duration = duration
	}
	s.EventCount++
	
	return s, nil
}

func (s *Session) MarkPersist() {
	s.Sign = 1
}

func (s *Session) MarkCollapse() {
	s.Sign = -1
}
