package main

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

func main() {

	router := newRouter()

	address := ":8080"

	log.Printf("ManaVault API listening on %s", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", notFoundHandler)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		response := errorResponse{
			Error: "method not allowed",
		}

		if err := writeJSON(w, http.StatusMethodNotAllowed, response); err != nil {
			log.Printf("Error encoding method not allowed response: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := healthResponse{Status: "ok"}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := errorResponse{Error: "route not found"}

	if err := writeJSON(w, http.StatusNotFound, response); err != nil {
		log.Printf("Error encoding not found response: %v", err)
	}
}
