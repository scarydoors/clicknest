package server

import (
	"log/slog"
	"net/http"

	"github.com/rs/cors"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/server/handlers"
	"github.com/scarydoors/clicknest/internal/stats"
)

func setupRoutes(mux *http.ServeMux, logger *slog.Logger, ingestService *ingest.Service, statsService *stats.Service) {
	mux.Handle("/", http.NotFoundHandler())

	cors := cors.AllowAll()

	apiMux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", cors.Handler(apiMux)))

	handlers.RegisterIngestRoutes(apiMux, logger, ingestService)
	handlers.RegisterStatsRoutes(apiMux, logger, statsService)
}

