package server

import (
	"encoding/json"
	"net/http"
)

// TODO: application/problem+json support
type apiError struct {

}

type handlerWithError interface {
	ServeHTTPWithError(http.ResponseWriter, *http.Request) error
}
type handlerWithErrorFunc func(http.ResponseWriter, *http.Request) error
func (h handlerWithErrorFunc) ServeHTTPWithError(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

func serveErrors(next handlerWithErrorFunc) http.Handler {
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
