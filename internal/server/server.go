package server

import (
	"log/slog"
	"net/http"
)

func NewServer(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, logger)

	return mux
}
