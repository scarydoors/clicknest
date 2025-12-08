package stats_test

import (
	"fmt"
	"testing"
	"time"
)

func mustParseDate(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParseDurationWithDate(t *testing.T) {
	tests := []struct {
		durationstr string
		startDate time.Time

		expected time.Duration
	}{
		{
			"2d",
			mustParseDate(time.DateOnly, "2020-12-04"),
			
			time.Duration(2 * 24 * time.Hour),
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.durationstr, tt.startDate)
		t.Run(testname, func (t *testing.T) {
			t.Errorf("fail!")
		})
	}
}
