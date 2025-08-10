package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/errorutil"
	"github.com/scarydoors/clicknest/internal/ingest"
)

func setupRoutes(mux *http.ServeMux, logger *slog.Logger, ingestService *ingest.Service) {
	mux.Handle("POST /event", handleEventPost(ingestService, logger))
	mux.Handle("/", http.NotFoundHandler())
}

type EventRequest struct {
	Domain string `json:"domain"`
	Kind   string `json:"kind"`
	Url    string `json:"url"`
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
				time.Now(),
				eventRequest.Domain,
				eventRequest.Kind,
				eventRequest.Url,
			)

			var salt uint64 = 0 // TODO
			ip := r.Header.Get("X-Forwarded-For")
			logger.Info("generating userid", "domain", event.Domain, "ip", ip, "ua", r.UserAgent())
			event.UserID = analytics.NewUserID(salt, event.Domain, ip, r.UserAgent())
			logger.Info("generated userid", "userid", event.UserID)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := ingestService.IngestEvent(ctx, event); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		},
	)
}
