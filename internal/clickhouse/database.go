package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickhouseDBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func ParseDSN(dsn string) (ClickhouseDBConfig, error) {
	u, err := url.Parse(dsn);
	if err != nil {
		return ClickhouseDBConfig{}, fmt.Errorf("clickhouse: parseDSN: %w", err)
	}

	password, set := u.User.Password();
	if !set {
		return ClickhouseDBConfig{}, errors.New("clickhouse: parseDSN: password is not set in DSN")
	}

	return ClickhouseDBConfig{
		Host: u.Hostname(),
		Port: u.Port(),
		Database: u.Path[1:],
		Username: u.User.Username(),
		Password: password,
	}, err
}

func NewClickhouseConn(ctx context.Context, config ClickhouseDBConfig) (driver.Conn, error) {
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

	if err := conn.Ping(ctx); err != nil {
		cerr := conn.Close()
		if cerr != nil {
			cerr = fmt.Errorf("conn close: %w", cerr)
		}
		return nil, errors.Join(fmt.Errorf("ping: %w", err), cerr)
	}

	return conn, nil
}
