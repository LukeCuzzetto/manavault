package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	address := ":8080"

	log.Printf("ManaVault API listening on %s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := healthResponse{Status: "ok"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
}
