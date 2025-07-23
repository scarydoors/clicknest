package clickhouse

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/scarydoors/clicknest/internal/event"
)

type EventRepository struct {
	conn driver.Conn
}

type EventModel struct {
	Timestamp time.Time `ch:"timestamp"`
	Domain string `ch:"domain"`
	Kind string `ch:"kind"`
	Pathname string `ch:"pathname"`
}

func (c *EventRepository) InsertEvent(ctx context.Context, e event.Event) error {
	err := c.conn.Exec(
		ctx,
		`INSERT INTO events (
			timestamp,
			domain,
			kind,
			pathname
		) VALUES (?, ?, ?, ?)`,
		e.Timestamp,
		e.Domain,
		e.Kind,
		e.Pathname,
	)

	if err != nil {
		return err
	}

	return nil;
}

func (c *EventRepository) AsyncInsertEvent(ctx context.Context, e event.Event) error {
	err := c.conn.AsyncInsert(
		ctx,
		`INSERT INTO events (
			timestamp,
			domain,
			kind,
			pathname
		) VALUES (?, ?, ?, ?)`,
		false,
		e.Timestamp,
		e.Domain,
		e.Kind,
		e.Pathname,
	)

	if err != nil {
		return err
	}

	return nil;
}

func (c *EventRepository) BatchInsertEvent(ctx context.Context, e []event.Event) error {
	batch, err := c.conn.PrepareBatch(ctx,
		`INSERT INTO events (
			timestamp,
			domain,
			kind,
			pathname
		)`,
	)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, evt := range e {
		err := batch.Append(
			evt.Timestamp,
			evt.Domain,
			evt.Kind,
			evt.Pathname,
		)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
