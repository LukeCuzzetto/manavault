package main

import (
	"log"
	"net/http"

	"github.com/LukeCuzzetto/manavault/internal/api"
	"github.com/LukeCuzzetto/manavault/internal/config"
)

func main() {
	cfg := config.Load()
	router := api.NewRouter()

	log.Printf("ManaVault API listening on %s", cfg.Address)

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
