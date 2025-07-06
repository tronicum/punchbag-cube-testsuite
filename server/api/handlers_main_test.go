package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/username/punchbag-cube-testsuite/server/models"
	"github.com/username/punchbag-cube-testsuite/server/store"
	"go.uber.org/zap"
)

// MockStore is a mock implementation of the Store interface
type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateCluster(cluster *models.Cluster) error {
	args := m.Called(cluster)
	return args.Error(0)
}

func (m *MockStore) GetCluster(id string) (*models.Cluster, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Cluster), args.Error(1)
}

func (m *MockStore) ListClusters() ([]*models.Cluster, error) {
	args := m.Called()
	return args.Get(0).([]*models.Cluster), args.Error(1)
}

func (m *MockStore) ListClustersByProvider(provider models.CloudProvider) ([]*models.Cluster, error) {
	args := m.Called(provider)
	return args.Get(0).([]*models.Cluster), args.Error(1)
}

func (m *MockStore) UpdateCluster(id string, cluster *models.Cluster) error {
	args := m.Called(id, cluster)
	return args.Error(0)
}

func (m *MockStore) DeleteCluster(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateTestResult(result *models.TestResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockStore) GetTestResult(id string) (*models.TestResult, error) {
	args := m.Called(id)
	return args.Get(0).(*models.TestResult), args.Error(1)
}

func (m *MockStore) ListTestResults(clusterID string) ([]*models.TestResult, error) {
	args := m.Called(clusterID)
	return args.Get(0).([]*models.TestResult), args.Error(1)
}

func (m *MockStore) UpdateTestResult(id string, result *models.TestResult) error {
	args := m.Called(id, result)
	return args.Error(0)
}

func (m *MockStore) CreateNodePool(nodePool *models.NodePool) error {
	args := m.Called(nodePool)
	return args.Error(0)
}

func (m *MockStore) GetNodePool(id string) (*models.NodePool, error) {
	args := m.Called(id)
	return args.Get(0).(*models.NodePool), args.Error(1)
}

func (m *MockStore) ListNodePools(clusterID string) ([]*models.NodePool, error) {
	args := m.Called(clusterID)
	return args.Get(0).([]*models.NodePool), args.Error(1)
}

func (m *MockStore) UpdateNodePool(id string, nodePool *models.NodePool) error {
	args := m.Called(id, nodePool)
	return args.Error(0)
}

func (m *MockStore) DeleteNodePool(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureMonitoring(monitoring *models.AzureMonitoring) error {
	args := m.Called(monitoring)
	return args.Error(0)
}

func (m *MockStore) GetAzureMonitoring(id string) (*models.AzureMonitoring, error) {
	args := m.Called(id)
	return args.Get(0).(*models.AzureMonitoring), args.Error(1)
}

func (m *MockStore) ListAzureMonitorings() ([]*models.AzureMonitoring, error) {
	args := m.Called()
	return args.Get(0).([]*models.AzureMonitoring), args.Error(1)
}

func (m *MockStore) UpdateAzureMonitoring(id string, monitoring *models.AzureMonitoring) error {
	args := m.Called(id, monitoring)
	return args.Error(0)
}

func (m *MockStore) DeleteAzureMonitoring(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureKubernetes(kubernetes *models.AzureKubernetes) error {
	args := m.Called(kubernetes)
	return args.Error(0)
}

func (m *MockStore) GetAzureKubernetes(id string) (*models.AzureKubernetes, error) {
	args := m.Called(id)
	return args.Get(0).(*models.AzureKubernetes), args.Error(1)
}

func (m *MockStore) ListAzureKubernetes() ([]*models.AzureKubernetes, error) {
	args := m.Called()
	return args.Get(0).([]*models.AzureKubernetes), args.Error(1)
}

func (m *MockStore) UpdateAzureKubernetes(id string, kubernetes *models.AzureKubernetes) error {
	args := m.Called(id, kubernetes)
	return args.Error(0)
}

func (m *MockStore) DeleteAzureKubernetes(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureBudget(budget *models.AzureBudget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockStore) GetAzureBudget(id string) (*models.AzureBudget, error) {
	args := m.Called(id)
	return args.Get(0).(*models.AzureBudget), args.Error(1)
}

func (m *MockStore) ListAzureBudgets() ([]*models.AzureBudget, error) {
	args := m.Called()
	return args.Get(0).([]*models.AzureBudget), args.Error(1)
}

func (m *MockStore) UpdateAzureBudget(id string, budget *models.AzureBudget) error {
	args := m.Called(id, budget)
	return args.Error(0)
}

func (m *MockStore) DeleteAzureBudget(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test setup helper
func setupTest() (*MockStore, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockStore := &MockStore{}
	router := gin.New()

	// Create handler with logger
	logger, _ := zap.NewDevelopment()

	// Setup routes using the existing SetupRoutes function
	SetupRoutes(router, mockStore, logger)

	return mockStore, router
}

// Test RunTest
func TestRunTest(t *testing.T) {
	mockStore, router := setupTest()

	t.Run("successful test run", func(t *testing.T) {
		// Mock that cluster exists
		cluster := &models.Cluster{
			ID:       "test-cluster-1",
			Name:     "Test Cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusRunning,
		}
		mockStore.On("GetCluster", "test-cluster-1").Return(cluster, nil)
		mockStore.On("CreateTestResult", mock.AnythingOfType("*models.TestResult")).Return(nil)

		testRequest := map[string]interface{}{
			"testType": "load-test",
		}

		body, _ := json.Marshal(testRequest)
		req, _ := http.NewRequest("POST", "/api/v1/clusters/test-cluster-1/tests", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		mockStore.On("GetCluster", "nonexistent").Return((*models.Cluster)(nil), store.ErrNotFound)

		testRequest := map[string]interface{}{
			"testType": "load-test",
		}

		body, _ := json.Marshal(testRequest)
		req, _ := http.NewRequest("POST", "/api/v1/clusters/nonexistent/tests", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/clusters/test-cluster-1/tests", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

// Test ListTestResults
func TestListTestResults(t *testing.T) {
	mockStore, router := setupTest()

	t.Run("successful test results listing", func(t *testing.T) {
		testResults := []*models.TestResult{
			{
				ID:        "test-1",
				ClusterID: "cluster-1",
				Status:    models.TestStatusPassed,
				TestType:  "load-test",
			},
			{
				ID:        "test-2",
				ClusterID: "cluster-1",
				Status:    models.TestStatusFailed,
				TestType:  "stress-test",
			},
		}

		mockStore.On("ListTestResults", "cluster-1").Return(testResults, nil)

		req, _ := http.NewRequest("GET", "/api/v1/clusters/cluster-1/tests", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response []*models.TestResult
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "test-1", response[0].ID)
		assert.Equal(t, "test-2", response[1].ID)

		mockStore.AssertExpectations(t)
	})

	t.Run("empty test results list", func(t *testing.T) {
		mockStore.On("ListTestResults", "cluster-1").Return([]*models.TestResult{}, nil)

		req, _ := http.NewRequest("GET", "/api/v1/clusters/cluster-1/tests", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response []*models.TestResult
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)

		mockStore.AssertExpectations(t)
	})
}

// Test ValidateProvider
func TestValidateProvider(t *testing.T) {
	mockStore, router := setupTest()

	t.Run("successful azure provider validation", func(t *testing.T) {
		validationRequest := map[string]interface{}{
			"provider": "azure",
			"credentials": map[string]interface{}{
				"subscription_id": "test-subscription",
				"tenant_id":       "test-tenant",
			},
		}

		body, _ := json.Marshal(validationRequest)
		req, _ := http.NewRequest("POST", "/api/v1/providers/validate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, true, response["valid"])
		assert.Equal(t, "azure", response["provider"])
	})

	t.Run("successful aws provider validation", func(t *testing.T) {
		validationRequest := map[string]interface{}{
			"provider": "aws",
			"credentials": map[string]interface{}{
				"access_key": "test-access-key",
				"secret_key": "test-secret-key",
				"region":     "us-west-2",
			},
		}

		body, _ := json.Marshal(validationRequest)
		req, _ := http.NewRequest("POST", "/api/v1/providers/validate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, true, response["valid"])
		assert.Equal(t, "aws", response["provider"])
	})

	t.Run("invalid provider", func(t *testing.T) {
		validationRequest := map[string]interface{}{
			"provider": "invalid-provider",
		}

		body, _ := json.Marshal(validationRequest)
		req, _ := http.NewRequest("POST", "/api/v1/providers/validate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/providers/validate", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

// Test SimulateProviderOperation
func TestSimulateProviderOperation(t *testing.T) {
	mockStore, router := setupTest()

	t.Run("successful azure operation simulation", func(t *testing.T) {
		simulationRequest := map[string]interface{}{
			"provider":  "azure",
			"operation": "create_cluster",
			"parameters": map[string]interface{}{
				"name":       "test-cluster",
				"node_count": 3,
				"vm_size":    "Standard_D2s_v3",
			},
		}

		body, _ := json.Marshal(simulationRequest)
		req, _ := http.NewRequest("POST", "/api/v1/providers/simulate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "azure", response["provider"])
		assert.Equal(t, "create_cluster", response["operation"])
		assert.NotNil(t, response["result"])
	})

	t.Run("successful aws operation simulation", func(t *testing.T) {
		simulationRequest := map[string]interface{}{
			"provider":  "aws",
			"operation": "create_cluster",
			"parameters": map[string]interface{}{
				"name":          "test-cluster",
				"instance_type": "t3.medium",
				"region":        "us-west-2",
			},
		}

		body, _ := json.Marshal(simulationRequest)
		req, _ := http.NewRequest("POST", "/api/v1/providers/simulate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "aws", response["provider"])
		assert.Equal(t, "create_cluster", response["operation"])
		assert.NotNil(t, response["result"])
	})

	t.Run("unsupported provider", func(t *testing.T) {
		simulationRequest := map[string]interface{}{
			"provider":  "unsupported",
			"operation": "create_cluster",
		}

		body, _ := json.Marshal(simulationRequest)
		req, _ := http.NewRequest("POST", "/api/v1/providers/simulate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/providers/simulate", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
