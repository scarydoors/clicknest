package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/scarydoors/clicknest/internal/analytics"
	"github.com/scarydoors/clicknest/internal/ingest"
)

func registerIngestRoutes(apiMux *http.ServeMux, logger *slog.Logger, ingestService *ingest.Service) {
	apiMux.Handle("POST /event", serveErrors(handleEventPost(ingestService, logger)))
}

type eventRequest struct {
	Domain string `json:"domain"`
	Kind   string `json:"kind"`
	Url    string `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	Data map[string]string `json:"data,omitempty"`
}

func handleEventPost(ingestService *ingest.Service, logger *slog.Logger) handlerWithErrorFunc {
	return handlerWithErrorFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			var eventRequest eventRequest
			if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
				return err
			}

			event, err := analytics.NewEvent(
				eventRequest.Timestamp,
				eventRequest.Domain,
				eventRequest.Kind,
				eventRequest.Url,
			)
			event.Data = eventRequest.Data

			var salt uint64 = 0 // TODO
			ip, err := getClientIP(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
			event.UserID = analytics.NewUserID(salt, event.Domain, ip, r.UserAgent())
			slog.Info("userid logged","salt", salt, "domain", event.Domain, "ip", ip, "ua", r.UserAgent())

			if err := ingestService.IngestEvent(ctx, event); err != nil {
				return err
			}
			w.WriteHeader(http.StatusNoContent)
			return nil
		},
	)
}
