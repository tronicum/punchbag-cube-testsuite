package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/username/punchbag-cube-testsuite/server/models"
	"github.com/username/punchbag-cube-testsuite/server/store"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockStore implements the Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateCluster(cluster *models.Cluster) error {
	args := m.Called(cluster)
	return args.Error(0)
}

func (m *MockStore) GetCluster(id string) (*models.Cluster, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cluster), args.Error(1)
}

func (m *MockStore) ListClusters() ([]*models.Cluster, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Cluster), args.Error(1)
}

func (m *MockStore) ListClustersByProvider(provider models.CloudProvider) ([]*models.Cluster, error) {
	args := m.Called(provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Cluster), args.Error(1)
}

func (m *MockStore) UpdateCluster(cluster *models.Cluster) error {
	args := m.Called(cluster)
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestResult), args.Error(1)
}

func (m *MockStore) ListTestResults() ([]*models.TestResult, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TestResult), args.Error(1)
}

func (m *MockStore) UpdateTestResult(result *models.TestResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockStore) CreateNodePool(nodePool *models.NodePool) error {
	args := m.Called(nodePool)
	return args.Error(0)
}

func (m *MockStore) GetNodePool(id string) (*models.NodePool, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NodePool), args.Error(1)
}

func (m *MockStore) ListNodePools() ([]*models.NodePool, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.NodePool), args.Error(1)
}

func (m *MockStore) UpdateNodePool(nodePool *models.NodePool) error {
	args := m.Called(nodePool)
	return args.Error(0)
}

func (m *MockStore) DeleteNodePool(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureBudget(budget *models.AzureBudget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockStore) GetAzureBudget(id string) (*models.AzureBudget, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AzureBudget), args.Error(1)
}

func (m *MockStore) ListAzureBudgets() ([]*models.AzureBudget, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

func (m *MockStore) CreateAzureMonitoring(monitoring *models.AzureMonitoring) error {
	args := m.Called(monitoring)
	return args.Error(0)
}

func (m *MockStore) GetAzureMonitoring(id string) (*models.AzureMonitoring, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AzureMonitoring), args.Error(1)
}

func (m *MockStore) ListAzureMonitorings() ([]*models.AzureMonitoring, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AzureKubernetes), args.Error(1)
}

func (m *MockStore) ListAzureKubernetes() ([]*models.AzureKubernetes, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

func setupTestRouter() (*gin.Engine, *MockStore, *Handlers) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockStore := &MockStore{}
	logger, _ := zap.NewDevelopment()
	handlers := NewHandlers(mockStore, logger)
	return router, mockStore, handlers
}

func TestCreateCluster(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.POST("/clusters", handlers.CreateCluster)

	t.Run("successful creation", func(t *testing.T) {
		cluster := models.Cluster{
			Name:     "test-cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusCreating,
		}

		mockStore.On("CreateCluster", mock.AnythingOfType("*models.Cluster")).Return(nil)

		body, _ := json.Marshal(cluster)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/clusters", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, cluster.Name, response.Name)
		assert.NotEmpty(t, response.ID) // ID should be generated
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster already exists", func(t *testing.T) {
		cluster := models.Cluster{
			Name:     "existing-cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusCreating,
		}

		mockStore.On("CreateCluster", mock.AnythingOfType("*models.Cluster")).Return(store.ErrAlreadyExists)

		body, _ := json.Marshal(cluster)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/clusters", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "cluster already exists")
		mockStore.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/clusters", bytes.NewBuffer([]byte("invalid json")))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetCluster(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.GET("/clusters/:id", handlers.GetCluster)

	t.Run("existing cluster", func(t *testing.T) {
		cluster := &models.Cluster{
			ID:        "cluster-1",
			Name:      "test-cluster",
			Provider:  models.Azure,
			Status:    models.ClusterStatusRunning,
			CreatedAt: time.Now(),
		}

		mockStore.On("GetCluster", "cluster-1").Return(cluster, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters/cluster-1", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, cluster.ID, response.ID)
		assert.Equal(t, cluster.Name, response.Name)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		mockStore.On("GetCluster", "nonexistent").Return(nil, store.ErrNotFound)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters/nonexistent", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "cluster not found")
		mockStore.AssertExpectations(t)
	})
}

func TestListClusters(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.GET("/clusters", handlers.ListClusters)

	t.Run("list all clusters", func(t *testing.T) {
		clusters := []*models.Cluster{
			{
				ID:       "cluster-1",
				Name:     "test-cluster-1",
				Provider: models.Azure,
				Status:   models.ClusterStatusRunning,
			},
			{
				ID:       "cluster-2",
				Name:     "test-cluster-2",
				Provider: models.AWS,
				Status:   models.ClusterStatusRunning,
			},
		}

		mockStore.On("ListClusters").Return(clusters, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		mockStore.AssertExpectations(t)
	})

	t.Run("list clusters by provider", func(t *testing.T) {
		azureClusters := []*models.Cluster{
			{
				ID:       "cluster-1",
				Name:     "azure-cluster",
				Provider: models.Azure,
				Status:   models.ClusterStatusRunning,
			},
		}

		mockStore.On("ListClustersByProvider", models.CloudProvider("azure")).Return(azureClusters, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters?provider=azure", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, models.Azure, response[0].Provider)
		mockStore.AssertExpectations(t)
	})
}

func TestUpdateCluster(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.PUT("/clusters/:id", handlers.UpdateCluster)

	t.Run("successful update", func(t *testing.T) {
		existingCluster := &models.Cluster{
			ID:       "cluster-1",
			Name:     "old-name",
			Provider: models.Azure,
			Status:   models.ClusterStatusRunning,
		}

		updateData := models.Cluster{
			Name:   "new-name",
			Status: "updating",
		}

		mockStore.On("GetCluster", "cluster-1").Return(existingCluster, nil)
		mockStore.On("UpdateCluster", mock.AnythingOfType("*models.Cluster")).Return(nil)

		body, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("PUT", "/clusters/cluster-1", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "new-name", response.Name)
		assert.Equal(t, "updating", response.Status)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		updateData := models.Cluster{Name: "new-name"}

		mockStore.On("GetCluster", "nonexistent").Return(nil, store.ErrNotFound)

		body, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("PUT", "/clusters/nonexistent", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockStore.AssertExpectations(t)
	})
}

func TestDeleteCluster(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.DELETE("/clusters/:id", handlers.DeleteCluster)

	t.Run("successful deletion", func(t *testing.T) {
		mockStore.On("DeleteCluster", "cluster-1").Return(nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("DELETE", "/clusters/cluster-1", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		mockStore.On("DeleteCluster", "nonexistent").Return(store.ErrNotFound)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("DELETE", "/clusters/nonexistent", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockStore.AssertExpectations(t)
	})
}

func TestRunTest(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.POST("/test", handlers.RunTest)

	t.Run("successful test run", func(t *testing.T) {
		cluster := &models.Cluster{
			ID:       "cluster-1",
			Name:     "test-cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusRunning,
		}

		testRequest := models.TestRequest{
			ClusterID: "cluster-1",
			TestType:  "connectivity",
			Config:    map[string]interface{}{"timeout": "30s"},
		}

		mockStore.On("GetCluster", "cluster-1").Return(cluster, nil)
		mockStore.On("CreateTestResult", mock.AnythingOfType("*models.TestResult")).Return(nil)

		body, _ := json.Marshal(testRequest)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var response models.TestResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testRequest.ClusterID, response.ClusterID)
		assert.Equal(t, testRequest.TestType, response.TestType)
		assert.NotEmpty(t, response.ID)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		testRequest := models.TestRequest{
			ClusterID: "nonexistent",
			TestType:  "connectivity",
		}

		mockStore.On("GetCluster", "nonexistent").Return(nil, store.ErrNotFound)

		body, _ := json.Marshal(testRequest)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "cluster not found")
		mockStore.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		testRequest := models.TestRequest{
			// Missing required fields
		}

		body, _ := json.Marshal(testRequest)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTestResult(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.GET("/test/:id", handlers.GetTestResult)

	t.Run("existing test result", func(t *testing.T) {
		testResult := &models.TestResult{
			ID:        "test-1",
			ClusterID: "cluster-1",
			TestType:  "connectivity",
			Status:    "completed",
			Duration:  30 * time.Second,
			StartedAt: time.Now().Add(-1 * time.Minute),
		}

		mockStore.On("GetTestResult", "test-1").Return(testResult, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test/test-1", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.TestResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testResult.ID, response.ID)
		assert.Equal(t, testResult.Status, response.Status)
		mockStore.AssertExpectations(t)
	})

	t.Run("test result not found", func(t *testing.T) {
		mockStore.On("GetTestResult", "nonexistent").Return(nil, store.ErrNotFound)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test/nonexistent", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockStore.AssertExpectations(t)
	})
}

func TestListTestResults(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.GET("/test", handlers.ListTestResults)

	t.Run("successful list", func(t *testing.T) {
		testResults := []*models.TestResult{
			{
				ID:        "test-1",
				ClusterID: "cluster-1",
				TestType:  "connectivity",
				Status:    "completed",
			},
			{
				ID:        "test-2",
				ClusterID: "cluster-2",
				TestType:  "performance",
				Status:    models.TestStatusPassed,
			},
		}

		mockStore.On("ListTestResults").Return(testResults, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.TestResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		mockStore.AssertExpectations(t)
	})
}

// Test provider validation from main handlers
func TestHandlersValidateProvider(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.GET("/validate/:provider", handlers.ValidateProvider)

	t.Run("azure provider validation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/validate/azure", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "azure")
		assert.Contains(t, w.Body.String(), "valid")
	})

	t.Run("invalid provider", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/validate/invalid", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "unsupported provider")
	})
}

// Test simulation from main handlers
func TestHandlersSimulateProviderOperation(t *testing.T) {
	router, mockStore, handlers := setupTestRouter()
	router.POST("/simulate", handlers.SimulateProviderOperation)

	t.Run("successful simulation", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"provider":  models.Azure,
			"operation": "create-cluster",
			"params": map[string]interface{}{
				"name": "test-cluster",
			},
		}

		body, _ := json.Marshal(requestBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/simulate", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "simulation")
	})
}

// Test additional edge cases and error scenarios
func TestHandlersErrorScenarios(t *testing.T) {
	t.Run("create cluster with existing ID", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.POST("/clusters", handlers.CreateCluster)

		cluster := &models.Cluster{
			ID:       "existing-cluster",
			Name:     "Test Cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusCreating,
		}

		mockStore.On("CreateCluster", mock.AnythingOfType("*models.Cluster")).Return(store.ErrAlreadyExists)

		body, _ := json.Marshal(cluster)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/clusters", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("list clusters with provider filter", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.GET("/clusters", handlers.ListClusters)

		azureClusters := []*models.Cluster{
			{
				ID:       "azure-cluster-1",
				Name:     "Azure Cluster 1",
				Provider: models.Azure,
				Status:   models.ClusterStatusRunning,
			},
		}

		mockStore.On("ListClustersByProvider", models.Azure).Return(azureClusters, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters?provider=azure", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, models.Azure, response[0].Provider)

		mockStore.AssertExpectations(t)
	})

	t.Run("run test on non-running cluster", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.POST("/clusters/:id/tests", handlers.RunTest)

		cluster := &models.Cluster{
			ID:       "creating-cluster",
			Name:     "Creating Cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusCreating,
		}

		mockStore.On("GetCluster", "creating-cluster").Return(cluster, nil)

		testRequest := map[string]interface{}{
			"testType": "load-test",
		}

		body, _ := json.Marshal(testRequest)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/clusters/creating-cluster/tests", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "cluster is not in running state")
		mockStore.AssertExpectations(t)
	})

	t.Run("validate provider with missing credentials", func(t *testing.T) {
		router, _, handlers := setupTestRouter()
		router.POST("/validate", handlers.ValidateProvider)

		requestBody := map[string]interface{}{
			"provider": models.Azure,
			// Missing credentials
		}

		body, _ := json.Marshal(requestBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/validate", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing required")
	})

	t.Run("simulate operation with invalid parameters", func(t *testing.T) {
		router, _, handlers := setupTestRouter()
		router.POST("/simulate", handlers.SimulateProviderOperation)

		requestBody := map[string]interface{}{
			"provider":  models.Azure,
			"operation": "invalid-operation",
		}

		body, _ := json.Marshal(requestBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/simulate", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Test concurrent operations
func TestConcurrentOperations(t *testing.T) {
	t.Run("concurrent cluster creation", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.POST("/clusters", handlers.CreateCluster)

		// Mock expects multiple cluster creations
		mockStore.On("CreateCluster", mock.AnythingOfType("*models.Cluster")).Return(nil).Times(5)

		// Create 5 concurrent requests
		results := make(chan int, 5)

		for i := 0; i < 5; i++ {
			go func(index int) {
				cluster := &models.Cluster{
					ID:       fmt.Sprintf("cluster-%d", index),
					Name:     fmt.Sprintf("Test Cluster %d", index),
					Provider: models.Azure,
					Status:   models.ClusterStatusCreating,
				}

				body, _ := json.Marshal(cluster)
				w := httptest.NewRecorder()
				r, _ := http.NewRequest("POST", "/clusters", bytes.NewBuffer(body))
				r.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, r)
				results <- w.Code
			}(i)
		}

		// Collect all results
		for i := 0; i < 5; i++ {
			statusCode := <-results
			assert.Equal(t, http.StatusCreated, statusCode)
		}

		mockStore.AssertExpectations(t)
	})

	t.Run("concurrent test runs", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.POST("/clusters/:id/tests", handlers.RunTest)

		cluster := &models.Cluster{
			ID:       "test-cluster",
			Name:     "Test Cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusRunning,
		}

		// Mock cluster exists and test results can be created
		mockStore.On("GetCluster", "test-cluster").Return(cluster, nil).Times(3)
		mockStore.On("CreateTestResult", mock.AnythingOfType("*models.TestResult")).Return(nil).Times(3)

		results := make(chan int, 3)

		for i := 0; i < 3; i++ {
			go func(index int) {
				testRequest := map[string]interface{}{
					"testType": fmt.Sprintf("test-%d", index),
				}

				body, _ := json.Marshal(testRequest)
				w := httptest.NewRecorder()
				r, _ := http.NewRequest("POST", "/clusters/test-cluster/tests", bytes.NewBuffer(body))
				r.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, r)
				results <- w.Code
			}(i)
		}

		// Collect all results
		for i := 0; i < 3; i++ {
			statusCode := <-results
			assert.Equal(t, http.StatusCreated, statusCode)
		}

		mockStore.AssertExpectations(t)
	})
}

// Test performance and load scenarios
func TestPerformanceScenarios(t *testing.T) {
	t.Run("large cluster list", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.GET("/clusters", handlers.ListClusters)

		// Create a large list of clusters
		clusters := make([]*models.Cluster, 1000)
		for i := 0; i < 1000; i++ {
			clusters[i] = &models.Cluster{
				ID:       fmt.Sprintf("cluster-%d", i),
				Name:     fmt.Sprintf("Cluster %d", i),
				Provider: models.Azure,
				Status:   models.ClusterStatusRunning,
			}
		}

		mockStore.On("ListClusters").Return(clusters, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1000)

		mockStore.AssertExpectations(t)
	})

	t.Run("large test result list", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.GET("/clusters/:id/tests", handlers.ListTestResults)

		// Create a large list of test results
		testResults := make([]*models.TestResult, 500)
		for i := 0; i < 500; i++ {
			testResults[i] = &models.TestResult{
				ID:        fmt.Sprintf("test-%d", i),
				ClusterID: "cluster-1",
				Status:    models.TestStatusPassed,
				TestType:  fmt.Sprintf("test-type-%d", i%5),
			}
		}

		mockStore.On("ListTestResults", "cluster-1").Return(testResults, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters/cluster-1/tests", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.TestResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 500)

		mockStore.AssertExpectations(t)
	})
}
