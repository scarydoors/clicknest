package clickhouse

import (
	"context"

	"github.com/scarydoors/clicknest/internal/event"
)

func (c *ClickhouseDB) InsertEvent(ctx context.Context, e event.Event) error {
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

func (c *ClickhouseDB) AsyncInsertEvent(ctx context.Context, e event.Event) error {
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

func (c *ClickhouseDB) BatchInsertEvent(ctx context.Context, e []event.Event) error {
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
