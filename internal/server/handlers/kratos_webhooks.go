package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/serverutil"
)

func RegisterKratosWebhooksRoutes(kratosMux *http.ServeMux, logger *slog.Logger) {
	kratosMux.Handle("POST /create-user", serverutil.ServeErrors(handleCreateUserPost(logger)))
}

func handleCreateUserPost(logger *slog.Logger) serverutil.HandlerWithErrorFunc {
	return serverutil.HandlerWithErrorFunc(
		func (w http.ResponseWriter, r *http.Request) error {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}

			logger.Info("received create-user", slog.String("body", string(body)))
			return fmt.Errorf("test");
		},
	)
}
