package sim

import (
	"encoding/json"
	"log"
	"net/http"
	mock "punchbag-cube-testsuite/multitool/pkg/mock"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// AKS Handlers
func HandleAks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var c sharedmodels.Cluster
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
		var ws sharedmodels.LogAnalyticsWorkspace
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
		var b sharedmodels.AzureBudget
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
		var a sharedmodels.AppInsightsResource
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

// AWS S3 Handler
func HandleS3(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var b sharedmodels.ObjectStorageBucket
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode S3 payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateS3(&b)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetS3(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: S3 not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListS3())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteS3(id) {
			log.Printf("[ERROR] %s %s: failed to delete S3 (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Azure Blob Storage Handler
func HandleBlobStorage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var b sharedmodels.ObjectStorageBucket
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode Blob payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateBlobStorage(&b)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetBlobStorage(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: Blob not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListBlobStorage())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteBlobStorage(id) {
			log.Printf("[ERROR] %s %s: failed to delete Blob (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Re-export unified handlers for sim-server
var HandleBlob = HandleBlobStorage

// GCS Handler
func HandleGCS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var b sharedmodels.ObjectStorageBucket
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			log.Printf("[ERROR] %s %s: failed to decode GCS payload: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateGCS(&b)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetGCS(id)
			if result == nil {
				log.Printf("[ERROR] %s %s: GCS not found (id=%s)", r.Method, r.URL.Path, id)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListGCS())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteGCS(id) {
			log.Printf("[ERROR] %s %s: failed to delete GCS (id=%s)", r.Method, r.URL.Path, id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
