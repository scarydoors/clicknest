package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/scarydoors/clicknest/internal/serverutil"
	"github.com/scarydoors/clicknest/internal/stats"
	//"github.com/go-playground/validator/v10"
)

func RegisterStatsRoutes(apiMux *http.ServeMux, logger *slog.Logger, statsService *stats.Service) {
	apiMux.Handle("GET /timeseries", serverutil.ServeErrors(handleTimeseriesGet(statsService, logger)))
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

// TODO: go-validator?
// type timeseriesGetRawParameters struct {
// 	startDate string
// 	endDate string
// 	interval string ``
// }
//
// type timeseriesGetParameters struct {
// 	startDate time.Time
// 	endDate time.Time
// 	interval time.Duration
// }

func handleTimeseriesGet(statsService *stats.Service, logger *slog.Logger) serverutil.HandlerWithErrorFunc {
	return serverutil.HandlerWithErrorFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			query := r.URL.Query()
			interval := query.Get("interval")
			if interval == "" {
				return errors.New("missing required query parameter: interval");
			}
			intervalDur, err := time.ParseDuration(interval)
			if err != nil {
				return errors.New("interval");
			}

			startDate := query.Get("start-date")
			if startDate == "" {
				return errors.New("missing required query parameter: interval");
			}
			startTime, err := time.Parse(time.RFC3339, startDate)
			if err != nil {
				return errors.New("missing required query parameter: start-date");
			}

			endDate := query.Get("end-date")
			if endDate == "" {
				return errors.New("missing required query parameter: end-date");
			}
			endTime, err := time.Parse(time.RFC3339, endDate)
			if err != nil {
				return errors.New("missing required query parameter: interval");
			}

			dur := endTime.Sub(startTime);
			estPoints := dur / intervalDur
			return errors.New(strconv.FormatInt(int64(estPoints), 10));

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
