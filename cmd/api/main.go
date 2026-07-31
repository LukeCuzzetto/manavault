package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LukeCuzzetto/manavault/internal/api"
	"github.com/LukeCuzzetto/manavault/internal/config"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverReadTimeout       = 10 * time.Second
	serverWriteTimeout      = 30 * time.Second
	serverIdleTimeout       = 60 * time.Second
)

func newHTTPServer(
	address string,
	handler http.Handler,
) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
	}
}

func main() {

	logger := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	app := api.NewApplication(logger)
	router := app.Router()

	server := newHTTPServer(cfg.Address, router)

	logger.Printf("ManaVault API listening on %s", cfg.Address)

	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("error running server: %v", err)
	}
}
