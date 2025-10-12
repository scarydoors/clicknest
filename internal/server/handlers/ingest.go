package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/ingest"
	"github.com/scarydoors/clicknest/internal/serverutil"
)

func RegisterIngestRoutes(apiMux *http.ServeMux, logger *slog.Logger, ingestService *ingest.Service) {
	apiMux.Handle("POST /event", handleEventPost(ingestService, logger))
}

type eventRequest struct {
	Domain string `json:"domain"`
	Kind   string `json:"kind"`
	Url    string `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	Data map[string]string `json:"data,omitempty"`
}

func handleEventPost(ingestService *ingest.Service, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var eventRequest eventRequest
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
			event.Data = eventRequest.Data

			var salt uint64 = 0 // TODO
			ip, err := serverutil.GetClientIP(r)
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
