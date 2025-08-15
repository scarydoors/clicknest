package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/errorutil"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/rs/cors"
)

func setupRoutes(mux *http.ServeMux, logger *slog.Logger, ingestService *ingest.Service) {
	mux.Handle("/", http.NotFoundHandler())

	corsMiddleware := cors.AllowAll()

	muxAPI := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", corsMiddleware.Handler(muxAPI)))

	muxAPI.Handle("/event", handleEventPost(ingestService, logger))
}

type EventRequest struct {
	Domain string `json:"domain"`
	Kind   string `json:"kind"`
	Url    string `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

func handleEventPost(ingestService *ingest.Service, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer errorutil.DeferIgnoreErr(r.Body.Close)

			ctx := r.Context()

			var eventRequest EventRequest
			if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			event, err := analytics.NewEvent(
				eventRequest.Timestamp,
				eventRequest.Domain,
				eventRequest.Kind,
				eventRequest.Url,
			)

			var salt uint64 = 0 // TODO
			ip, err := getClientIP(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			event.UserID = analytics.NewUserID(salt, event.Domain, ip, r.UserAgent())
			slog.Info("userid logged","salt", salt, "domain", event.Domain, "ip", ip, "ua", r.UserAgent())

			if err := ingestService.IngestEvent(ctx, event); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	)
}
