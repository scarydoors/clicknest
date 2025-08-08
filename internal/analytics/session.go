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
	start      time.Time
	end        time.Time
	domain     string
	duration   uint32
	eventCount uint32
	sessionID  SessionID
	userID     UserID
	sign       int8
}

func FromEvent(event Event) Session {
	return Session{
		start:      event.Timestamp,
		end:        event.Timestamp,
		domain:     event.Domain,
		duration:   0,
		eventCount: 1,
		sessionID:  NewSessionID(),
		userID:     event.UserID,
		sign:       1,
	}
}

func (s *Session) Update(event Event) {
	s.end = event.Timestamp
	s.duration = uint32(s.end.Sub(s.start).Milliseconds()) 
}
