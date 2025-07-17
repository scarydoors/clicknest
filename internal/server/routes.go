package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/event"
)

func setupRoutes(mux *http.ServeMux, logger *slog.Logger) {
	mux.Handle("POST /event", handleEventPost(logger))
	mux.Handle("/", http.NotFoundHandler())
}

func handleEventPost(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			var event event.Event

			if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			logger.Info("recieved event", slog.Any("event", event))
		},
	)
}
