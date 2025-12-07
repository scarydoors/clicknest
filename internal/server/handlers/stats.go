package handlers

import (
	"encoding/json"
	"fmt"

	//"errors"
	"log/slog"
	"net/http"

	//"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/scarydoors/clicknest/internal/serverutil"
	"github.com/scarydoors/clicknest/internal/stats"
	"github.com/gorilla/schema"
	//"github.com/go-playground/locales/en"
	//ut "github.com/go-playground/universal-translator"
	//en_translations "github.com/go-playground/validator/v10/translations/en"
)

func RegisterStatsRoutes(apiMux *http.ServeMux, logger *slog.Logger, validate *validator.Validate, statsService *stats.Service) {
	apiMux.Handle("GET /timeseries", serverutil.ServeErrors(handleTimeseriesGet(statsService, logger, validate)))
}

type timeseriesGetRawParameters struct {
	StartDate string `schema:"start_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate string `schema:"end_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Interval string `schema:"interval" validate:"required,duration"`
}

func (t timeseriesGetRawParameters) ToParams() (stats.GetTimeseriesParameters, error) {
	startDate, err := time.Parse(time.RFC3339, t.StartDate)
	if err != nil {
		return stats.GetTimeseriesParameters{}, err
	}

	endDate, err := time.Parse(time.RFC3339, t.EndDate)
	if err != nil {
		return stats.GetTimeseriesParameters{}, err
	}

	interval, err := time.ParseDuration(t.Interval)
	if err != nil {
		return stats.GetTimeseriesParameters{}, err
	}

	return stats.GetTimeseriesParameters{
		StartDate: startDate,
		EndDate: endDate,
		Interval: interval,
	}, nil
}

func handleTimeseriesGet(statsService *stats.Service, logger *slog.Logger, validate *validator.Validate) serverutil.HandlerWithErrorFunc {
	return serverutil.HandlerWithErrorFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			decoder := schema.NewDecoder()

			var rawParams timeseriesGetRawParameters
			if err := decoder.Decode(&rawParams, r.URL.Query()); err != nil {
				return err
			}

			if err := validate.Struct(rawParams); err != nil {
				return err
			}

			params, err := rawParams.ToParams()
			if err != nil {
				return err
			}

			timeseries, err := statsService.GetTimeseries(ctx, params);
			fmt.Printf("%+v\n", timeseries)
			if err != nil {
				return err
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(timeseries)

			return nil;
		},
	)
}
