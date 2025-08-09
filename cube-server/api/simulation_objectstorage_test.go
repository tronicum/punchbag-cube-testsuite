package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
	store "github.com/tronicum/punchbag-cube-testsuite/store"
	"go.uber.org/zap"
)

// NewTestSimulationService returns a SimulationService for tests
func NewTestSimulationService() *simulation.SimulationService {
	return simulation.NewSimulationService()
}

func TestSimulateS3ObjectStorageOperations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Provide a mock store, logger, and simulation service
	var mockStore store.Store = nil
	logger := zap.NewNop()
	sim := NewTestSimulationService()
	SetupRoutes(r, mockStore, logger, sim)

	// Create bucket
	createReq := map[string]interface{}{
		"provider":  "aws",
		"operation": "create_bucket",
		"parameters": map[string]interface{}{
			"name":   "test-bucket-sim",
			"region": "eu-central-1",
		},
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/simulate/providers/aws/operations/create_bucket", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("create_bucket: expected 200, got %d", resp.Code)
	}

	// Set bucket policy
	policyReq := map[string]interface{}{
		"provider":  "aws",
		"operation": "set_bucket_policy",
		"parameters": map[string]interface{}{
			"bucket": "test-bucket-sim",
			"policy": "{\"Version\":\"2012-10-17\",\"Statement\":[]}",
		},
	}
	body, _ = json.Marshal(policyReq)
	req = httptest.NewRequest("POST", "/api/v1/simulate/providers/aws/operations/set_bucket_policy", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("set_bucket_policy: expected 200, got %d", resp.Code)
	}

	// Set bucket versioning
	versioningReq := map[string]interface{}{
		"provider":  "aws",
		"operation": "set_bucket_versioning",
		"parameters": map[string]interface{}{
			"bucket":  "test-bucket-sim",
			"enabled": true,
		},
	}
	body, _ = json.Marshal(versioningReq)
	req = httptest.NewRequest("POST", "/api/v1/simulate/providers/aws/operations/set_bucket_versioning", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("set_bucket_versioning: expected 200, got %d", resp.Code)
	}

	// Set bucket lifecycle
	lifecycleReq := map[string]interface{}{
		"provider":  "aws",
		"operation": "set_bucket_lifecycle",
		"parameters": map[string]interface{}{
			"bucket":    "test-bucket-sim",
			"lifecycle": "{\"Rules\":[{\"ID\":\"expire-objects\",\"Status\":\"Enabled\",\"Expiration\":{\"Days\":1}}]}",
		},
	}
	body, _ = json.Marshal(lifecycleReq)
	req = httptest.NewRequest("POST", "/api/v1/simulate/providers/aws/operations/set_bucket_lifecycle", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("set_bucket_lifecycle: expected 200, got %d", resp.Code)
	}
}
