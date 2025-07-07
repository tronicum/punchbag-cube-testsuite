package sim

import (
	"encoding/json"
	"net/http"
	mock "punchbag-cube-testsuite/multitool/pkg/mock"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// AWS EKS Handlers
func HandleEks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var c sharedmodels.Cluster
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateEks(&c)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetEks(id)
			if result == nil { w.WriteHeader(http.StatusNotFound); return }
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListEks())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteEks(id) { w.WriteHeader(http.StatusNotFound); return }
		w.WriteHeader(http.StatusNoContent)
	}
}

// GCP GKE Handlers
func HandleGke(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var c sharedmodels.Cluster
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// TODO: Implement GKE mock logic
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodGet, http.MethodDelete:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
