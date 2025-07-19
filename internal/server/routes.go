package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/scarydoors/clicknest/internal/event"
)

func setupRoutes(mux *http.ServeMux, logger *slog.Logger, eventService *event.Service) {
	mux.Handle("POST /event", handleEventPost(eventService, logger))
	mux.Handle("/", http.NotFoundHandler())
}

type EventRequest struct {
	Domain string `json:"domain"`
	Kind   string `json:"kind"`
	Url    string `json:"url"`
}

func handleEventPost(eventService *event.Service, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			ctx := r.Context()

			var eventRequest EventRequest
			if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			event, err := event.NewEvent(
				time.Now(),
				eventRequest.Domain,
				eventRequest.Kind,
				eventRequest.Url,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := eventService.IngestEvent(ctx, event); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		},
	)
}
