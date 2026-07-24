package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", notFoundHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("Error encoding method not allowed response: %v", err)
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
		log.Printf("Error encoding health response: %v", err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := errorResponse{
		Error: "route not found",
	}

	if err := writeJSON(
		w,
		http.StatusNotFound,
		response,
	); err != nil {
		log.Printf("Error encoding not found response: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
