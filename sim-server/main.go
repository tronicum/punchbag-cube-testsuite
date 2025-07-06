package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"punchbag-cube-testsuite/sim"
)

func main() {
	logFile, err := os.OpenFile("sim-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)

	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Cube Simulation Server is running...")
	})

	http.HandleFunc("/api/simulate/azure/aks", sim.HandleAks)
	http.HandleFunc("/api/simulate/azure/loganalytics", sim.HandleLogAnalytics)
	http.HandleFunc("/api/simulate/azure/budget", sim.HandleBudget)
	http.HandleFunc("/api/simulate/azure/appinsights", sim.HandleAppInsights)

	http.HandleFunc("/api/simulate/aws/eks", sim.HandleEks)
	http.HandleFunc("/api/simulate/aws/s3", sim.HandleS3)
	http.HandleFunc("/api/simulate/gcp/gke", sim.HandleGke)

	http.HandleFunc("/api/validation", sim.HandleValidation)

	// Mock JWT and OAuth endpoints
	http.HandleFunc("/api/mock-jwt-login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token": "mock-jwt-token", "type": "jwt"}`))
	})
	http.HandleFunc("/api/mock-oauth-login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token": "mock-oauth-token", "type": "oauth"}`))
	})

	// Simple API key auth middleware
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "dev-secret-key"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/status" { // allow status without auth
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-API-Key") != apiKey {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: missing or invalid API key"))
			return
		}
		// fallback to default mux
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	fmt.Println("Cube Simulation Server is running on :8080...")
	http.ListenAndServe(":8080", nil)
}
