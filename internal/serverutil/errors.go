package serverutil

import (
	"encoding/json"
	"net/http"
)

// TODO: application/problem+json support
type apiError struct {

}

type HandlerWithError interface {
	ServeHTTPWithError(http.ResponseWriter, *http.Request) error
}
type HandlerWithErrorFunc func(http.ResponseWriter, *http.Request) error
func (h HandlerWithErrorFunc) ServeHTTPWithError(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

func ServeErrors(next HandlerWithErrorFunc) http.Handler {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			err := next.ServeHTTPWithError(w, r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"type": "cool-error",
					"error": err.Error(),
				});
			}
		},
	)
}
