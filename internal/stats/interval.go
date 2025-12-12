package stats

import (
	"errors"
	"strconv"
	"unicode"
)

type IntervalUnit int

const (
	IntervalUnitSecond IntervalUnit = iota
	IntervalUnitMinute
	IntervalUnitHour
	IntervalUnitDay
	IntervalUnitWeek
	IntervalUnitMonth
	IntervalUnitQuarter
	IntervalUnitYear
)

const maxUnitLen int = 2
var unitKey = map[string]IntervalUnit{
	"s": IntervalUnitSecond,
	"m": IntervalUnitMinute,
	"h": IntervalUnitHour,
	"d": IntervalUnitDay,
	"w": IntervalUnitWeek,
	"mo": IntervalUnitMonth,
	"q": IntervalUnitQuarter,
	"y": IntervalUnitYear,
}

type Interval struct {
	Value float64
	Unit IntervalUnit
}

var ErrInvalidInterval = errors.New("invalid interval")

func ParseInterval(s string) (Interval, error) {
	if s == "" {
		return Interval{}, ErrInvalidInterval
	}

	var splitIdx int = -1
	for i, r := range s {
		if unicode.IsLetter(r) {
			splitIdx = i
			break
		}
	}

	// if not unit is found or unit is at index 0
	if splitIdx <= 0 {
		return Interval{}, ErrInvalidInterval
	}

	valuePart := s[:splitIdx]
	unitPart := s[splitIdx:]

	value, err := strconv.ParseFloat(valuePart, 64)
	if err != nil {
		return Interval{}, ErrInvalidInterval
	}

	unit, ok := unitKey[unitPart]
	if !ok {
		return Interval{}, ErrInvalidInterval
	}

	return Interval{
		Value: value,
		Unit: unit,
	}, nil
}
