package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LukeCuzzetto/manavault/internal/api"
	"github.com/LukeCuzzetto/manavault/internal/config"
)

func main() {

	logger := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	app := api.NewApplication(logger)
	router := app.Router()

	logger.Printf("ManaVault API listening on %s", cfg.Address)

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		logger.Fatalf("error starting server: %v", err)
	}
}
