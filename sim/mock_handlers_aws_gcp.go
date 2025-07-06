package sim

import (
	"encoding/json"
	"net/http"
	"punchbag-cube-testsuite/multitool/pkg/mock"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

// AWS EKS Handlers
func HandleEks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var c models.Cluster
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
		var c models.Cluster
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateGke(&c)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetGke(id)
			if result == nil { w.WriteHeader(http.StatusNotFound); return }
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListGke())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteGke(id) { w.WriteHeader(http.StatusNotFound); return }
		w.WriteHeader(http.StatusNoContent)
	}
}

// AWS S3 Handlers
func HandleS3(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var b models.S3Bucket
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := mock.MockCreateS3(&b)
		json.NewEncoder(w).Encode(result)
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			result := mock.MockGetS3(id)
			if result == nil { w.WriteHeader(http.StatusNotFound); return }
			json.NewEncoder(w).Encode(result)
			return
		}
		json.NewEncoder(w).Encode(mock.MockListS3())
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" || !mock.MockDeleteS3(id) { w.WriteHeader(http.StatusNotFound); return }
		w.WriteHeader(http.StatusNoContent)
	}
}
