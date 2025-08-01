package clickhouse

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/scarydoors/clicknest/internal/analytics"
)

type EventRepository struct {
	conn driver.Conn
}

type EventModel struct {
	Timestamp time.Time `ch:"timestamp"`
	Domain string `ch:"domain"`
	Kind string `ch:"kind"`
	Pathname string `ch:"pathname"`
	UserId uint64 `ch:"user_id"`
}

func NewEventRepository(conn driver.Conn) *EventRepository {
	return &EventRepository{
		conn: conn,
	}
}

func marshalEvent(event analytics.Event) EventModel {
	return EventModel{
		Timestamp: event.Timestamp,
		Domain: event.Domain,
		Kind: event.Kind,
		Pathname: event.Pathname,
		UserId: event.UserId,
	}
}

func (c *EventRepository) BatchInsertEvent(ctx context.Context, events []analytics.Event) error {
	batch, err := c.conn.PrepareBatch(ctx,
		`INSERT INTO events (
			timestamp,
			domain,
			kind,
			pathname,
		    user_id
		)`,
	)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, event := range events {
		model := marshalEvent(event)	
		err := batch.AppendStruct(&model)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
