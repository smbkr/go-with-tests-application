package main

import (
	"application/src/data"
	"application/src/server"
	"log"
	"net/http"
)

func main() {
	store := data.NewInMemoryPlayerStore()
	handler := &server.PlayerServer{store}
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatalf("unable to listen on port :5000, %v", err)
	}
}
