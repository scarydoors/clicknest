package server

import (
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/stats"
)

func NewServer(logger *slog.Logger, ingestService *ingest.Service, statsService *stats.Service) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, logger, ingestService, statsService)

	return mux
}
