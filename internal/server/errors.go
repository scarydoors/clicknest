package server

import "net/http"

type apiError struct {

}

type handlerWithError func(http.ResponseWriter, *http.Request) error
