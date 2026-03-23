package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func registerKratosWebhooksRoutes(kratosMux *http.ServeMux, logger *slog.Logger) {
	kratosMux.Handle("POST /create-user", serveErrors(handleCreateUserPost(logger)))
}

type createUserPostParameters struct {
	UserId string `validate:"required,uuid" json:"user_id"`
	Email string `validate:"required,email" json:"email"`
	FirstName string `validate:"required" json:"first_name"`
	LastName string `validate:"required" json:"last_name"`
}

func handleCreateUserPost(logger *slog.Logger) handlerWithErrorFunc {
	return handlerWithErrorFunc(
		func (w http.ResponseWriter, r *http.Request) error {
			var params createUserPostParameters
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&params); err != nil {
				return err
			}

			logger.Info("handleCreateUserPost: received user info", slog.Any("params", params))

			w.WriteHeader(http.StatusNoContent)
			return nil
		},
	)
}
