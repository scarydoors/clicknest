package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func getPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", getPing)

	fmt.Printf("Serving server at port :6969\n")
	if err := http.ListenAndServe(":6969", mux); errors.Is(err, http.ErrServerClosed){
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s", err)
		os.Exit(1)
	}
}
