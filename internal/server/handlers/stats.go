package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
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

			query := r.URL.Query()
			interval := query.Get("interval")
			if interval == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "missing required query parameter: interval",
				})
				return;
			}
			intervalDur, err := time.ParseDuration(interval)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "interval",
				})
				return;
			}

			startDate := query.Get("start-date")
			if startDate == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "missing required query parameter: start-date",
				})
				return;
			}
			startTime, err := time.Parse(time.RFC3339, startDate)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "start-time",
				})
				return;
			}

			endDate := query.Get("end-date")
			if endDate == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "missing required query parameter: end-date",
				})
				return;
			}
			endTime, err := time.Parse(time.RFC3339, endDate)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "end-time",
				})
				return;
			}

			dur := endTime.Sub(startTime);
			estPoints := dur / intervalDur
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": strconv.FormatInt(int64(estPoints), 10),
			})
			return;

			timeseries, err := statsService.GetPageviews(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(timeseriesToTimeseriesResponse(timeseries))
		},
	)
}
