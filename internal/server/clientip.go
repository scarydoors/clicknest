package server

import "net/http"

const xffHeaderName = "X-Forwarded-For"

// TODO: fix naive implementation of IP handler
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get(xffHeaderName); xff != "" {
		return xff
	}

	return r.RemoteAddr
}
