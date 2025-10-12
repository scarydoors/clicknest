package clickhouse

import (
	"context"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/scarydoors/clicknest/internal/errorutil"
	"github.com/scarydoors/clicknest/internal/stats"
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

type timeseriesPoint struct {
	Timestamp time.Time `ch:"timestamp"`
	Value uint64 `ch:"value"`
}

type timeseries []timeseriesPoint
func (t timeseries) toPageviewResult() stats.Timeseries {
	res := make(stats.Timeseries, 0, len(t))
	for _, point := range t {
		res = append(res, stats.TimeseriesPoint{
			Timestamp: point.Timestamp,
			Value: point.Value,
		})
	}

	return res
}


func (s *StatsRepository) GetPageviews(ctx context.Context) (stats.Timeseries, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT
			toStartOfHour(timestamp) as timestamp,
			count() AS value
		FROM events
		WHERE kind = 'pageview'
		GROUP BY timestamp
		ORDER BY timestamp ASC
		WITH FILL STEP INTERVAL 1 HOUR
		`)
	if err != nil {
		return nil, err
	}
	defer errorutil.DeferErrf(&err, "rows close: %w", rows.Close)

	var timeseries timeseries
	for rows.Next() {
		var timeseriesPoint timeseriesPoint
		if err := rows.ScanStruct(&timeseriesPoint); err != nil {
			return nil, err
		}
		timeseries = append(timeseries, timeseriesPoint)
	}

	return timeseries.toPageviewResult(), rows.Err()
}
