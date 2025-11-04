package handlers

import (
	"encoding/json"

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

func RegisterStatsRoutes(apiMux *http.ServeMux, logger *slog.Logger, validate *validator.Validate, statsService *stats.Service) {
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
	Interval string `validate:"required,duration"`
}

type timeseriesGetParameters struct {
	StartDate time.Time
	EndDate time.Time
	Interval time.Duration `validate:"interval_granularity=StartDate~EndDate:1000"`
}

func timeseriesGetParamsFromRawParams(rawParams timeseriesGetRawParameters) (timeseriesGetParameters, error) {
	startDate, err := time.Parse(time.RFC3339, rawParams.StartDate)
	if err != nil {
		return timeseriesGetParameters{}, err
	}

	endDate, err := time.Parse(time.RFC3339, rawParams.EndDate)
	if err != nil {
		return timeseriesGetParameters{}, err
	}

	interval, err := time.ParseDuration(rawParams.Interval)
	if err != nil {
		return timeseriesGetParameters{}, err
	}

	return timeseriesGetParameters{
		StartDate: startDate,
		EndDate: endDate,
		Interval: interval,
	}, nil
}

func handleTimeseriesGet(statsService *stats.Service, logger *slog.Logger, validate *validator.Validate) serverutil.HandlerWithErrorFunc {
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

			err := validate.Struct(rawParams)
			if err != nil {
				return err
			}

			params, err := timeseriesGetParamsFromRawParams(rawParams)
			if err != nil {
				return err
			}

			err = validate.Struct(params)
			if err != nil {
				return err
			}
			
			return err


			timeseries, err := statsService.GetPageviews(ctx);
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(timeseriesToTimeseriesResponse(timeseries))

			return nil;
		},
	)
}
