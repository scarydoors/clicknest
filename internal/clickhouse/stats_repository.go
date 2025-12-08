package clickhouse

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
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


func (s *StatsRepository) GetPageviews(ctx context.Context, params stats.GetTimeseriesParameters) (stats.Timeseries, error) {
	intervalSeconds, err := DurationToIntervalSeconds(params.Interval)
	if err != nil {
		return stats.Timeseries{}, err
	}

	rows, err := s.conn.Query(ctx, `
		SELECT
			toStartOfInterval(timestamp, INTERVAL {interval: UInt64} SECONDS) as timestamp,
			count() AS value
		FROM events
		WHERE kind = 'pageview'
		AND timestamp >= {start_date: DateTime}
		AND timestamp <= {end_date: DateTime}
		AND domain = {domain: String}
		GROUP BY timestamp
		ORDER BY timestamp ASC
		WITH FILL STEP INTERVAL {interval: UInt64} SECONDS`,
		clickhouse.Named("start_date", params.StartDate.Format("2006-01-02 15:04:05")),
		clickhouse.Named("end_date", params.EndDate.Format("2006-01-02 15:04:05")),
		clickhouse.Named("interval", strconv.FormatUint(intervalSeconds, 10)),
		clickhouse.Named("domain", params.Domain),
	)
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
