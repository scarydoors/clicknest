package clickhouse

import (
	"fmt"
	"math"
	"time"
)

func DurationToIntervalSeconds(duration time.Duration) (uint64, error) {
	// the timestamps are stored as DateTime rather than DateTime64 which means
	// that the precision of the timestamps is less than a second
	seconds := math.Floor(duration.Abs().Seconds())
	
	if seconds > float64(math.MaxUint64) {
		return 0, fmt.Errorf("overflow converting %f into uint64", seconds)
	}

	return uint64(seconds), nil
}
