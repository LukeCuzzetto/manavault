package main

import (
	"log"
	"net/http"

	"github.com/LukeCuzzetto/manavault/internal/api"
)

func main() {
	router := api.NewRouter()

	address := ":8080"

	log.Printf("ManaVault API listening on %s", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
