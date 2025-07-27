package server

import (
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/ingest"
)

func NewServer(logger *slog.Logger, ingestService *ingest.Service) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, logger, ingestService)

	return mux
}
