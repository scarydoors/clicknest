package server

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/stats"
)

func NewServer(logger *slog.Logger, validate *validator.Validate, ingestService *ingest.Service, statsService *stats.Service) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, logger, validate, ingestService, statsService)

	return mux
}
