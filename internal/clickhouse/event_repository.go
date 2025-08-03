package clickhouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/errorutil"
)

type EventRepository struct {
	conn   driver.Conn
	logger *slog.Logger
}

type EventModel struct {
	Timestamp time.Time `ch:"timestamp"`
	Domain    string    `ch:"domain"`
	Kind      string    `ch:"kind"`
	SessionId uint64    `ch:"session_id"`
	UserId    uint64    `ch:"user_id"`
	Pathname  string    `ch:"pathname"`
}

func NewEventRepository(conn driver.Conn, logger *slog.Logger) *EventRepository {
	return &EventRepository{
		conn:   conn,
		logger: logger,
	}
}

func marshalEvent(event analytics.Event) EventModel {
	return EventModel{
		Timestamp: event.Timestamp,
		Domain:    event.Domain,
		Kind:      event.Kind,
		SessionId: event.SessionId,
		UserId:    event.UserId,
		Pathname:  event.Pathname,
	}
}

func (c *EventRepository) BatchInsert(ctx context.Context, events []analytics.Event) (err error) {
	batch, err := c.conn.PrepareBatch(ctx,
		`INSERT INTO events (
			timestamp,
			domain,
			kind,
			session_id,
		    user_id,
			pathname
		)`,
	)
	if err != nil {
		return err
	}
	defer errorutil.DeferErrf(&err, "batch close: %w", batch.Close)

	for _, event := range events {
		model := marshalEvent(event)
		err := batch.AppendStruct(&model)
		if err != nil {
			return err
		}
	}

	c.logger.Info("batch inserted events", slog.Int("count", len(events)))
	if err := batch.Send(); err != nil {
		return fmt.Errorf("batch send: %w", err)
	}

	return nil
}
