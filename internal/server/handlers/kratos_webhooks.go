package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/scarydoors/clicknest/internal/serverutil"
)

func RegisterKratosWebhooksRoutes(kratosMux *http.ServeMux, logger *slog.Logger) {
	kratosMux.Handle("POST /create-user", serverutil.ServeErrors(handleCreateUserPost(logger)))
}

type CreateUserPostParameters struct {
	UserId string `validate:"required,uuid" json:"user_id"`
	Email string `validate:"required,email" json:"email"`
	FirstName string `validate:"required" json:"first_name"`
	LastName string `validate:"required" json:"last_name"`
}

func handleCreateUserPost(logger *slog.Logger) serverutil.HandlerWithErrorFunc {
	return serverutil.HandlerWithErrorFunc(
		func (w http.ResponseWriter, r *http.Request) error {
			var params CreateUserPostParameters
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return err
			}

			logger.Info("handleCreateUserPost: received user info", slog.Any("params", params))

			return fmt.Errorf("test");
		},
	)
}
