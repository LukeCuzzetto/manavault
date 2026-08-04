package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LukeCuzzetto/deckengine/internal/api"
	"github.com/LukeCuzzetto/deckengine/internal/config"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverReadTimeout       = 10 * time.Second
	serverWriteTimeout      = 30 * time.Second
	serverIdleTimeout       = 60 * time.Second
	serverShutdownTimeout   = 10 * time.Second
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

	shutdownSignal, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	logger.Printf("DECK//ENGINE API is listening on %s", cfg.Address)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("error running server: %v", err)
		}

		return

	case <-shutdownSignal.Done():
		logger.Println("shutdown signal received, shutting down server...")
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		serverShutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Printf("graceful shut down failed: %v", err)
	} else {
		logger.Println("DECK//ENGINE API Stopped")
	}
}
