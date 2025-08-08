package analytics

import (
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
	Duration   uint32
	EventCount uint32
	SessionID  SessionID
	UserID     UserID
	Sign       int8
}

func FromEvent(event Event) Session {
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
