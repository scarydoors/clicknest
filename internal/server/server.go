package server

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	kratos "github.com/ory/kratos-client-go"
	"github.com/rs/cors"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/stats"
)

func NewServer(
	logger *slog.Logger,
	validate *validator.Validate,
	kratosClient *kratos.APIClient,
	ingestService *ingest.Service,
	statsService *stats.Service,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())

	cors := cors.AllowAll()

	apiMux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", cors.Handler(apiMux)))

	registerIngestRoutes(apiMux, logger, ingestService)
	registerStatsRoutes(apiMux, logger, validate, statsService)

	kratosMux := http.NewServeMux()
	mux.Handle("/kratos-webhooks/", http.StripPrefix("/kratos-webhooks", kratosMux))
	registerKratosWebhooksRoutes(kratosMux, logger)
	return mux
}
