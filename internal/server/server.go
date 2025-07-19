package server

import (
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/event"
)

func NewServer(logger *slog.Logger, eventService *event.Service) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, logger, eventService)

	return mux
}
