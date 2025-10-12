package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/scarydoors/clicknest/internal/stats"
)

func RegisterStatsRoutes(apiMux *http.ServeMux, logger *slog.Logger, statsService *stats.Service) {
	apiMux.Handle("GET /timeseries", handleTimeseriesGet(statsService, logger))
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

func handleTimeseriesGet(statsService *stats.Service, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			timeseries, err := statsService.GetPageviews(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(timeseriesToTimeseriesResponse(timeseries))
		},
	)
}
