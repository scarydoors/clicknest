package server

import (
	"net"
	"net/http"
)

const xffHeaderName = "X-Forwarded-For"

// TODO: fix naive implementation of IP handler
func getClientIP(r *http.Request) (string, error) {
	var ip string
	if xff := r.Header.Get(xffHeaderName); xff != "" {
		ip = xff
	} else {
		ip = r.RemoteAddr
	}

	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return "", err
	}

	return host, nil
}
