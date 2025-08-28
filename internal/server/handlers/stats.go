package handlers

import (
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/stats"
)

func RegisterStatsRoutes(apiMux *http.ServeMux, logger *slog.Logger, statsService *stats.Service) {
	apiMux.Handle("GET /graph", handleGraphGet(statsService, logger))
}

func handleGraphGet(statsService *stats.Service, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello World"))
		},
	)
}
