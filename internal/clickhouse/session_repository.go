package clickhouse

import (
	"context"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/scarydoors/clicknest/internal/analytics"
)

type SessionRepository struct {
	conn driver.Conn
	logger *slog.Logger
}

type SessionModel struct {
	Start time.Time `ch:"start"`
	End time.Time `ch:"end"`
	Domain string `ch:"domain"`
	Duration uint32 `ch:"duration"`
	EventCount uint32 `ch:"event_count"`
	SessionId uint64 `ch:"session_id"`
	UserId uint64 `ch:"user_id"`
	Sign int8 `ch:"sign"`
}

func NewSessionRepository(conn driver.Conn, logger *slog.Logger) *SessionRepository {
	return &SessionRepository{
		conn: conn,
		logger: logger,
	}
}

func marshalSession(session analytics.Session) SessionModel {
	return SessionModel{
		Start: session.Start,
		End: session.End,
		Domain: session.Domain,
		Duration: session.Duration,
		EventCount: session.EventCount,
		SessionId: session.SessionId,
		UserId: session.UserId,
		Sign: session.Sign,
	}
}

func (c *SessionRepository) BatchInsert(ctx context.Context, sessions []analytics.Session) error {
	batch, err := c.conn.PrepareBatch(ctx,
		`INSERT INTO sessions (
			start,
			end,
			domain,
			duration,
			event_count,
			session_id,
			user_id,
			sign
		)`,
	)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, session := range sessions {
		model := marshalSession(session)	
		err := batch.AppendStruct(&model)
		if err != nil {
			return err
		}
	}

	c.logger.Info("batch inserted sessions", slog.Int("count", len(sessions)))
	return batch.Send()
}
