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
