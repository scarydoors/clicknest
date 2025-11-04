package validatorutil

import (
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const intervalGranularityValidatorTag string = "interval_granularity"

type intervalGranularityParam struct {
	minFieldName string
	maxFieldName string
	maxGranularity int
}

func parseIntervalGranularityParam(param string) (intervalGranularityParam, bool) {
	parts := strings.Split(param, ":")
	if len(parts) != 2 {
		return intervalGranularityParam{}, false
	}

	fieldsPart := parts[0]
	granularityPart := parts[1]

	maxGranularity, err := strconv.Atoi(granularityPart)	
	if err != nil {
		return intervalGranularityParam{}, false
	}

	fields := strings.Split(fieldsPart, "~")
	if len(fields) != 2 {
		return intervalGranularityParam{}, false
	}
	minFieldName := strings.TrimSpace(fields[0])
	maxFieldName := strings.TrimSpace(fields[1])

	return intervalGranularityParam{
		minFieldName: minFieldName,
		maxFieldName: maxFieldName,
		maxGranularity: maxGranularity,
	}, true
}

func makeIntervalGranularityValidator(cache *paramCache, logger *slog.Logger) func(validator.FieldLevel) bool {
	logger = logger.With(slog.Group("validator", slog.String("tag", intervalGranularityValidatorTag)))
	return func (fl validator.FieldLevel) bool {
		param := fl.Param()
		logger := logger.With(slog.Group("validator", slog.String("param", param)))

		parsedParam, ok := getParam(cache, intervalGranularityValidatorTag, param, parseIntervalGranularityParam)
		if !ok {
			logger.Error("unable to parse param for validator")
			return false
		}

		parent := fl.Parent()
		minField := parent.FieldByName(parsedParam.maxFieldName)
		maxField := parent.FieldByName(parsedParam.minFieldName)

		if !minField.IsValid() || !maxField.IsValid() {
			return false
		}

		timeType := reflect.TypeOf(time.Time{}) 
		if minField.Type() != timeType || maxField.Type() != timeType {
			return false
		}
		minTime := minField.Interface().(time.Time)
		maxTime := maxField.Interface().(time.Time)

		field := fl.Field()
		if field.Kind() != reflect.Int64 {
			return false
		}
		interval := time.Duration(field.Int())

		timeRange := maxTime.Sub(minTime)

		granules := timeRange / interval
		if int(granules) > parsedParam.maxGranularity {
			return false
		}


		return true
	}
}
