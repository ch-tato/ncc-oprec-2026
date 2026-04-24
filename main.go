package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var startTime time.Time

type HealthResponse struct {
	Name      string `json:"name"`
	NRP       string `json:"nrp"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Name:      "Muhammad Quthbi Danish Abqori",
		NRP:       "5025241036",
		Status:    "Sukses Sarimi Isi 200 OK",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	startTime = time.Now()

	port := 1234

	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/health", healthHandler)

	fmt.Printf("Server is running on port %s...\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
