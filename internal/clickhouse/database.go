package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickhouseDBConfig struct {
	Host string
	Port string
	Database string
	Username string
	Password string
}

func NewClickhouseConn(context context.Context, config ClickhouseDBConfig) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("open conn: %w", err)
	}

	if err := conn.Ping(context); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	return conn, nil
}
