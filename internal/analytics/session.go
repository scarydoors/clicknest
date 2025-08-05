package analytics

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

type Session struct {
	Start      time.Time
	End        time.Time
	Domain     string
	Duration   uint32
	EventCount uint32
	SessionID  uint64
	UserID     uint64
	Sign       int8
}

func FromEvent(event Event) Session {
	return Session{
		Start:      event.Timestamp,
		End:        event.Timestamp,
		Domain:     event.Domain,
		Duration:   0,
		EventCount: 1,
		SessionID:  generateSessionID(),
		UserID:     event.UserID,
		Sign:       1,
	}
}

func generateSessionID() uint64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(b[:])
}
