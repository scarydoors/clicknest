package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	Domain string `json:"domain"`
	Type string `json:"type"`
	Url string `json:"url"`
}

func handleEventPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var event Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("recieved event: %+v\n", event)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /event", handleEventPost)

	fmt.Printf("Serving server at port :6969\n")
	if err := http.ListenAndServe(":6969", mux); errors.Is(err, http.ErrServerClosed){
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s", err)
		os.Exit(1)
	}
}
