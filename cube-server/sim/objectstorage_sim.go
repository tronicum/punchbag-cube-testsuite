package sim

import (
	   "encoding/json"
	   "fmt"
	   "net/http"
	   "os"
	   "strings"
	   "sync"

	   sharederrors "github.com/tronicum/punchbag-cube-testsuite/shared/errors"
	   "github.com/tronicum/punchbag-cube-testsuite/shared/log"
)

// ObjectStorageSim defines a generic interface for object storage simulation
type ObjectStorageSim interface {
	CreateBucket(name string) error
	PutObject(bucket, key string, value []byte) error
	GetObject(bucket, key string) ([]byte, error)
	DeleteBucket(name string) error
	ListBuckets() []string
	ListObjects(bucket string) []string
}

// S3Object represents a simple S3 object (key, value).
type S3Object struct {
	Key   string
	Value []byte
}

// HetznerS3Mock simulates a minimal Hetzner S3 API in-memory.
type HetznerS3Mock struct {
	buckets map[string]map[string]S3Object // bucket -> key -> object
	mu      sync.RWMutex
}

// NewHetznerS3Mock creates a new in-memory Hetzner S3 mock.
func NewHetznerS3Mock() *HetznerS3Mock {
	return &HetznerS3Mock{
		buckets: make(map[string]map[string]S3Object),
	}
}

// DeleteBucket removes a bucket and all its objects.
func (m *HetznerS3Mock) DeleteBucket(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.buckets[name]; !exists {
		return sharederrors.ErrNotFound
	}
	delete(m.buckets, name)
	return nil
}

// ObjectStorageSimHandler exposes any ObjectStorageSim as an HTTP handler for simulation
func ObjectStorageSimHandler(sim ObjectStorageSim) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		   path := r.URL.Path
		   w.Header().Set("Content-Type", "application/json")
		   // Debug log for all incoming simulation requests
		   fmt.Printf("[SIM DEBUG] %s %s\n", r.Method, path)
		   // Accept both '/buckets' and '/buckets/' (trailing slash)
		   if strings.HasSuffix(path, "/buckets") || strings.HasSuffix(path, "/buckets/") {
			   switch r.Method {
			   case http.MethodGet:
				   // List buckets
				   buckets := sim.ListBuckets()
				   type ObjectStorageBucket struct {
					   ID             string                 `json:"id,omitempty"`
					   Name           string                 `json:"name"`
					   Provider       string                 `json:"provider"`
					   Region         string                 `json:"region,omitempty"`
					   Location       string                 `json:"location,omitempty"`
					   CreatedAt      string                 `json:"created_at,omitempty"`
					   UpdatedAt      string                 `json:"updated_at,omitempty"`
					   Policy         interface{}            `json:"policy,omitempty"`
					   Lifecycle      interface{}            `json:"lifecycle,omitempty"`
					   ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
				   }
				   var out []ObjectStorageBucket
				   for _, name := range buckets {
					   out = append(out, ObjectStorageBucket{
						   ID: "sim-" + name,
						   Name: name,
						   Provider: "hetzner",
						   Region: "fsn1",
						   CreatedAt: "2025-01-01T00:00:00Z",
						   UpdatedAt: "2025-01-01T00:00:00Z",
						   Policy: nil,
						   Lifecycle: nil,
						   ProviderConfig: map[string]interface{}{},
					   })
				   }
				   w.WriteHeader(http.StatusOK)
				   _ = json.NewEncoder(w).Encode(out)
			   case http.MethodPost:
				   // Create bucket
				   var req struct {
					   Name           string                 `json:"name"`
					   Provider       string                 `json:"provider"`
					   Region         string                 `json:"region"`
					   Location       string                 `json:"location,omitempty"`
					   Policy         interface{}            `json:"policy,omitempty"`
					   Lifecycle      interface{}            `json:"lifecycle,omitempty"`
					   ProviderConfig map[string]interface{} `json:"provider_config,omitempty"`
				   }
				   if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					   w.WriteHeader(http.StatusBadRequest)
					   _ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
					   return
				   }
				   err := sim.CreateBucket(req.Name)
				   if err != nil {
					   if err == sharederrors.ErrConflict {
						   w.WriteHeader(http.StatusConflict)
						   _ = json.NewEncoder(w).Encode(map[string]string{"error": "Bucket already exists"})
					   } else {
						   w.WriteHeader(http.StatusInternalServerError)
						   _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
					   }
					   return
				   }
				   w.WriteHeader(http.StatusCreated)
				   _ = json.NewEncoder(w).Encode(map[string]interface{}{
					   "id": "sim-" + req.Name,
					   "name": req.Name,
					   "provider": req.Provider,
					   "region": req.Region,
					   "location": req.Location,
					   "created_at": "2025-01-01T00:00:00Z",
					   "updated_at": "2025-01-01T00:00:00Z",
					   "policy": req.Policy,
					   "lifecycle": req.Lifecycle,
					   "provider_config": req.ProviderConfig,
				   })
			   default:
				   w.WriteHeader(http.StatusNotImplemented)
				   _ = json.NewEncoder(w).Encode(map[string]string{"error": "Not implemented for this method"})
			   }
			   return
		   }
		   w.WriteHeader(http.StatusNotFound)
		   _ = json.NewEncoder(w).Encode(map[string]string{"error": "Unknown endpoint"})
   })
}

// ListBuckets returns a list of all bucket names.
func (m *HetznerS3Mock) ListBuckets() []string {
	   // If SIMULATE_DUMMY_S3_CREDS=1, always return a dummy bucket for simulation/testing
	   if os.Getenv("SIMULATE_DUMMY_S3_CREDS") == "1" {
			   return []string{"sim-bucket-1"}
	   }
	   m.mu.RLock()
	   defer m.mu.RUnlock()
	   var names []string
	   for name := range m.buckets {
			   names = append(names, name)
	   }
	   return names
}

// ListObjects returns a list of all object keys in a bucket.
func (m *HetznerS3Mock) ListObjects(bucket string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var keys []string
	b, ok := m.buckets[bucket]
	if !ok {
		return keys
	}
	for k := range b {
		keys = append(keys, k)
	}
	return keys
}

// CreateBucket creates a new bucket.
func (m *HetznerS3Mock) CreateBucket(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.buckets[name]; exists {
		return sharederrors.ErrConflict
	}
	m.buckets[name] = make(map[string]S3Object)
	return nil
}

// PutObject puts an object into a bucket.
func (m *HetznerS3Mock) PutObject(bucket, key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return sharederrors.ErrNotFound
	}
	b[key] = S3Object{Key: key, Value: value}
	return nil
}

// GetObject retrieves an object from a bucket.
func (m *HetznerS3Mock) GetObject(bucket, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	obj, ok := b[key]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return obj.Value, nil
}

// HetznerS3MockHandler exposes the HetznerS3Mock as an HTTP handler for simulation
func HetznerS3MockHandler(mock *HetznerS3Mock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Warn("Hetzner S3 simulation endpoint called: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte("Hetzner S3 simulation not yet implemented"))
	})
}

// StartObjectStorageSimulationServer starts a simulation server for S3-like APIs
func StartObjectStorageSimulationServer(provider, port string) {
	switch provider {
	case "hetzner":
		fmt.Printf("Starting Hetzner S3 simulation on http://localhost:%s\n", port)
		mock := NewHetznerS3Mock()
		http.Handle("/", HetznerS3MockHandler(mock))
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Provider not supported for simulation.")
		os.Exit(1)
	}
}
