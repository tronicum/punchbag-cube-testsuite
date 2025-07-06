package sim

import (
	"encoding/json"
	"log"
	"net/http"
	"punchbag-cube-testsuite/multitool/pkg/mock"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

// AKS Handlers
func HandleAks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var c models.Cluster
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode AKS payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateAks(&c)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetAks(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: AKS not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListAks())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteAks(id) {
			log.Printf("[ERROR] %s %s: failed to delete AKS (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Log Analytics Handlers
func HandleLogAnalytics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var ws models.LogAnalyticsWorkspace
		if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode Log Analytics payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateLogAnalytics(&ws)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetLogAnalytics(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: Log Analytics not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListLogAnalytics())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteLogAnalytics(id) {
			log.Printf("[ERROR] %s %s: failed to delete Log Analytics (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Budget Handlers
func HandleBudget(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var b models.AzureBudget
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode Budget payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateBudget(&b)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetBudget(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: Budget not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListBudgets())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteBudget(id) {
			log.Printf("[ERROR] %s %s: failed to delete Budget (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// App Insights Handlers
func HandleAppInsights(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var a models.AppInsightsResource
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode App Insights payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateAppInsights(&a)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetAppInsights(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: App Insights not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListAppInsights())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteAppInsights(id) {
			log.Printf("[ERROR] %s %s: failed to delete App Insights (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
