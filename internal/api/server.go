package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Application struct {
	logger *log.Logger
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
		app.logger.Printf(
			"request method=%s path%s",
			r.Method,
			r.URL.Path,
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
