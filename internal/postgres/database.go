package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)


func NewPostgresConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("open conn: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		cerr := conn.Close(ctx)
		if cerr != nil {
			cerr = fmt.Errorf("conn close: %w", cerr)
		}
		return nil, errors.Join(fmt.Errorf("ping: %w", err), cerr)
	}

	return conn, nil
}
