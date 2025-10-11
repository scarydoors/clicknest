package clickhouse

import (
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type StatsRepository struct {
	conn driver.Conn	
	logger *slog.Logger
}

func NewStatsRepository(conn driver.Conn, logger *slog.Logger) *StatsRepository {
	return &StatsRepository{
		conn: conn,
		logger: logger,
	}
}
