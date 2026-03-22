package server

import (
	"net/http"

	kratos "github.com/ory/kratos-client-go"
	"github.com/scarydoors/clicknest/internal/auth"
)

func newRequireAuthMiddleware(kratosClient *kratos.APIClient) func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func (w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				r.Cookie(auth.SessionCookieName)
				kratosClient.FrontendAPI.ToSession(ctx).Cookie()
			},
		)
	}
}
