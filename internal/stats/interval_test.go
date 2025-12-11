package stats_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/scarydoors/clicknest/internal/stats"
)

func mustParseDate(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParseIntervalWithDate(t *testing.T) {
	tests := []struct {
		intervalstr string
		expected stats.Interval
	}{
		{
			"2d",
			stats.Interval{2, stats.IntervalUnitDay},
		},
		{
			"1.50m",
			stats.Interval{1.5, stats.IntervalUnitMinute},
		},
		{
			"1.508.65d",
			stats.Interval{1.5, stats.IntervalUnitMinute},
		},
		{
			"150865",
			stats.Interval{1.5, stats.IntervalUnitMinute},
		},
		{
			"d",
			stats.Interval{1.5, stats.IntervalUnitMinute},
		},
		{
			"56",
			stats.Interval{1.5, stats.IntervalUnitMinute},
		},
		{
			"8.e56d789",
			stats.Interval{1.5, stats.IntervalUnitMinute},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.intervalstr)
		t.Run(testname, func (t *testing.T) {
			interval, err := stats.ParseInterval(tt.intervalstr)
			t.Fatalf("%+v %+v", interval, err)
		})
	}
}
