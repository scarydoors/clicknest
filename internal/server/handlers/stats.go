package handlers

import (
	"encoding/json"
	"reflect"
	//"errors"
	"log/slog"
	"net/http"

	//"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/scarydoors/clicknest/internal/serverutil"
	"github.com/scarydoors/clicknest/internal/stats"
	//"github.com/go-playground/locales/en"
	//ut "github.com/go-playground/universal-translator"
	//en_translations "github.com/go-playground/validator/v10/translations/en"
)

func RegisterStatsRoutes(apiMux *http.ServeMux, logger *slog.Logger, statsService *stats.Service) {
	//en := en.New()
	//uni := ut.New(en, en)

	//trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	apiMux.Handle("GET /timeseries", serverutil.ServeErrors(handleTimeseriesGet(statsService, logger, validate)))
}

type timeseriesResponsePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value uint64 `json:"value"`
}

type timeseriesResponse []timeseriesResponsePoint

func timeseriesToTimeseriesResponse(ts stats.Timeseries) timeseriesResponse {
	timeseriesResp := make(timeseriesResponse, 0, len(ts))
	for _, t := range ts {
		timeseriesResp = append(timeseriesResp, timeseriesResponsePoint{
			Timestamp: t.Timestamp,
			Value: t.Value,
		})
	}

	return timeseriesResp
}

type timeseriesGetRawParameters struct {
	StartDate string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate string `validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Interval string
	Parsed timeseriesGetParameters `validate:"-"`
}

func timeseriesGetRawParametersValidation(sl validator.StructLevel) {
	rawParams := sl.Current().Interface().(timeseriesGetRawParameters)

	startTime, err := time.Parse(time.RFC3339, rawParams.StartDate)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(time.RFC3339, rawParams.EndDate)
	if err != nil {
		panic(err)
	}

	// TODO: custom duration validation
	interval, err := time.ParseDuration(rawParams.Interval)
	if err != nil {
		panic(err)
	}

	dur := endTime.Sub(startTime);
	estPoints := dur / interval

	if estPoints > 1000 {
		sl.ReportError(rawParams.Interval, "interval", "interval", "intervalgranularity", "")
	}

	p := sl.Current().FieldByName("Parsed")
	if p.IsValid() && p.CanSet() {
		p.Set(reflect.ValueOf(timeseriesGetParameters{
				startDate: startTime,
				endDate: endTime,
				interval: interval,
			}))
	}
}

type timeseriesGetParameters struct {
	startDate time.Time
	endDate time.Time
	interval time.Duration
}

func handleTimeseriesGet(statsService *stats.Service, logger *slog.Logger, validator *validator.Validate) serverutil.HandlerWithErrorFunc {
	validator.RegisterStructValidation(timeseriesGetRawParametersValidation, timeseriesGetRawParameters{});

	return serverutil.HandlerWithErrorFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			query := r.URL.Query()
			startDate := query.Get("start-date")
			endDate := query.Get("end-date")
			interval := query.Get("interval")

			rawParams := timeseriesGetRawParameters{
				StartDate: startDate,
				EndDate: endDate,
				Interval: interval,
			}

			err := validator.Struct(&rawParams)
			logger.Info("rawParams", slog.Any("struct", rawParams), slog.Any("yeah", err))
			return err


			timeseries, err := statsService.GetPageviews(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(timeseriesToTimeseriesResponse(timeseries))

			return nil;
		},
	)
}
