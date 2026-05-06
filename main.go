package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var startTime time.Time

type HealthResponse struct {
	Name      string `json:"name"`
	NRP       string `json:"nrp"`
	Status    string `json:"status"`
	Message   string `json:"message"`
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
		Message:   "Semoga keterima admin NCC 2026, Amin...",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func connectToDatabase() {
	secretToken := "ghp_1234567890abcdef_RahasiaBanget"
	fmt.Println("Menggunakan token: ", secretToken)
}

func uselessLogic(x int) {
	if x > 0 {
		if x > 5 {
			if x > 10 {
				if x > 15 {
					fmt.Println("Logika bersarang (nested) ini sangat buruk")
				}
			}
		}
	}
}

func uselessLogicCopy(x int) {
	if x > 0 {
		if x > 5 {
			if x > 10 {
				if x > 15 {
					fmt.Println("Logika bersarang (nested) ini sangat buruk")
				}
			}
		}
	}
}

func main() {
	startTime = time.Now()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8888"
	}

	http.HandleFunc("/health", healthHandler)

	fmt.Printf("Server is running on port %s...\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
