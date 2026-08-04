package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const databaseReadinessTimeout = 2 * time.Second

type Application struct {
	logger   *log.Logger
	database Database
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}
type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type readinessResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func NewApplication(logger *log.Logger, database Database) *Application {
	return &Application{
		logger:   logger,
		database: database,
	}
}

func (app *Application) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.healthHandler)
	mux.HandleFunc("/ready", app.readinessHandler)
	mux.HandleFunc("/", app.notFoundHandler)

	return app.requestLogger(app.recoverPanic(mux))
}

func (app *Application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		recorder := &responseRecorder{
			ResponseWriter: w,
		}

		next.ServeHTTP(recorder, r)

		app.logger.Printf(
			"request method=%s path=%s status=%d duration=%s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			time.Since(start),
		)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				app.logger.Printf("panic recovered: %v", recovered)

				response := errorResponse{
					Error: "internal server error",
				}

				if err := writeJSON(
					w,
					http.StatusInternalServerError,
					response,
				); err != nil {
					app.logger.Printf(
						"error encoding panic response: %v",
						err,
					)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		response := errorResponse{
			Error: "method not allowed",
		}

		if err := writeJSON(
			w,
			http.StatusMethodNotAllowed,
			response,
		); err != nil {
			app.logger.Printf("Error encoding method not allowed response: %v", err)
		}
		return
	}

	response := healthResponse{
		Status: "ok",
	}

	if err := writeJSON(
		w,
		http.StatusOK,
		response,
	); err != nil {
		app.logger.Printf("Error encoding health response: %v", err)
	}
}

func (app *Application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := errorResponse{
		Error: "route not found",
	}

	if err := writeJSON(
		w,
		http.StatusNotFound,
		response,
	); err != nil {
		app.logger.Printf("Error encoding not found response: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func (recorder *responseRecorder) WriteHeader(statusCode int) {
	if recorder.statusCode != 0 {
		return
	}

	recorder.statusCode = statusCode
	recorder.ResponseWriter.WriteHeader(statusCode)
}

func (recorder *responseRecorder) Write(data []byte) (int, error) {
	if recorder.statusCode == 0 {
		recorder.WriteHeader(http.StatusOK)
	}

	return recorder.ResponseWriter.Write(data)
}

func (app *Application) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		response := errorResponse{
			Error: "method not allowed",
		}

		if err := writeJSON(
			w,
			http.StatusMethodNotAllowed,
			response,
		); err != nil {
			app.logger.Printf("Error encoding method not allowed response: %v", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), databaseReadinessTimeout)
	defer cancel()

	if err := app.database.Ping(ctx); err != nil {
		app.logger.Printf("Database readiness check failed: %v", err)

		response := errorResponse{
			Error: "service unavailable",
		}

		if err := writeJSON(
			w,
			http.StatusServiceUnavailable,
			response,
		); err != nil {
			app.logger.Printf("Error encoding readiness error response: %v", err)
		}
		return

	}

	response := readinessResponse{
		Status:   "ready",
		Database: "ok",
	}

	if err := writeJSON(
		w,
		http.StatusOK,
		response,
	); err != nil {
		app.logger.Printf("Error encoding readiness response: %v", err)
	}
}
