package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Application struct {
	logger *log.Logger
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

func NewApplication(logger *log.Logger) *Application {
	return &Application{
		logger: logger,
	}
}

func (app *Application) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.healthHandler)
	mux.HandleFunc("/", app.notFoundHandler)

	return app.requestLogger(mux)
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
		recorder.statusCode = http.StatusOK
	}

	return recorder.ResponseWriter.Write(data)
}
