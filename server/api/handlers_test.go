package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MockStore struct{}

func (m *MockStore) SaveObjectStorageBucket(bucket *sharedmodels.ObjectStorageBucket) error {
	return nil
}
func (m *MockStore) GetObjectStorageBucket(provider, name string) (*sharedmodels.ObjectStorageBucket, error) {
	return &sharedmodels.ObjectStorageBucket{
		Provider: sharedmodels.CloudProvider(provider),
		Name:    name,
		Region:  "us-east-1",
	}, nil
}

// Implement all methods required by the store.Store interface as no-ops for testing
func (m *MockStore) CreateCluster(cluster *sharedmodels.Cluster) error { return nil }
func (m *MockStore) GetCluster(id string) (*sharedmodels.Cluster, error) { return nil, nil }
func (m *MockStore) ListClusters() ([]*sharedmodels.Cluster, error) { return nil, nil }
func (m *MockStore) ListClustersByProvider(provider sharedmodels.CloudProvider) ([]*sharedmodels.Cluster, error) { return nil, nil }
func (m *MockStore) UpdateCluster(id string, cluster *sharedmodels.Cluster) error { return nil }
func (m *MockStore) DeleteCluster(id string) error { return nil }
func (m *MockStore) CreateTestResult(result *sharedmodels.TestResult) error { return nil }
func (m *MockStore) GetTestResult(id string) (*sharedmodels.TestResult, error) { return nil, nil }
func (m *MockStore) ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error) { return nil, nil }
func (m *MockStore) UpdateTestResult(id string, result *sharedmodels.TestResult) error { return nil }
func (m *MockStore) CreateNodePool(nodePool *sharedmodels.NodePool) error { return nil }
func (m *MockStore) GetNodePool(id string) (*sharedmodels.NodePool, error) { return nil, nil }
func (m *MockStore) ListNodePools(clusterID string) ([]*sharedmodels.NodePool, error) { return nil, nil }
func (m *MockStore) UpdateNodePool(id string, nodePool *sharedmodels.NodePool) error { return nil }
func (m *MockStore) DeleteNodePool(id string) error { return nil }
func (m *MockStore) CreateAzureMonitoring(monitoring *sharedmodels.AzureMonitoring) error { return nil }
func (m *MockStore) GetAzureMonitoring(id string) (*sharedmodels.AzureMonitoring, error) { return nil, nil }
func (m *MockStore) ListAzureMonitorings() ([]*sharedmodels.AzureMonitoring, error) { return nil, nil }
func (m *MockStore) UpdateAzureMonitoring(id string, monitoring *sharedmodels.AzureMonitoring) error { return nil }
func (m *MockStore) DeleteAzureMonitoring(id string) error { return nil }
func (m *MockStore) CreateAzureKubernetes(kubernetes *sharedmodels.AzureKubernetes) error { return nil }
func (m *MockStore) GetAzureKubernetes(id string) (*sharedmodels.AzureKubernetes, error) { return nil, nil }
func (m *MockStore) ListAzureKubernetes() ([]*sharedmodels.AzureKubernetes, error) { return nil, nil }
func (m *MockStore) UpdateAzureKubernetes(id string, kubernetes *sharedmodels.AzureKubernetes) error { return nil }
func (m *MockStore) DeleteAzureKubernetes(id string) error { return nil }
func (m *MockStore) CreateAzureBudget(budget *sharedmodels.AzureBudget) error { return nil }
func (m *MockStore) GetAzureBudget(id string) (*sharedmodels.AzureBudget, error) { return nil, nil }
func (m *MockStore) ListAzureBudgets() ([]*sharedmodels.AzureBudget, error) { return nil, nil }
func (m *MockStore) UpdateAzureBudget(id string, budget *sharedmodels.AzureBudget) error { return nil }
func (m *MockStore) DeleteAzureBudget(id string) error { return nil }

func TestProxyObjectStorage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := &MockStore{}
	logger := zap.NewNop()
	h := NewHandlers(mockStore, logger)

	bucket := sharedmodels.ObjectStorageBucket{
		Provider: sharedmodels.CloudProvider("aws"),
		Name:     "test-bucket",
		Region:   "us-east-1",
	}
	body, _ := json.Marshal(bucket)
	req := httptest.NewRequest("POST", "/api/proxy/aws/objectstorage", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ProxyObjectStorage(c)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	var got sharedmodels.ObjectStorageBucket
	if err := json.Unmarshal(respBody, &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.Name != bucket.Name || got.Provider != bucket.Provider {
		t.Errorf("unexpected response: %+v", got)
	}
}
