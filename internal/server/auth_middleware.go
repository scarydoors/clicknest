package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	kratos "github.com/ory/kratos-client-go"
)

func newRequireAuthMiddleware(kratosClient *kratos.APIClient, logger *slog.Logger) func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func (w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()


				cookie := r.Header.Get("cookie")
				session, resp, err := kratosClient.FrontendAPI.ToSession(ctx).Cookie(cookie).Execute()
				if err != nil {
					fmt.Printf("%v\n",resp);
					// TODO: respond json function
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{
						"type": "cool-error",
						"error": err.Error(),
					});
					return;
				}

				next.ServeHTTP(w, r.WithContext(contextSetSession(ctx, session)))
			},
		)
	}
}
