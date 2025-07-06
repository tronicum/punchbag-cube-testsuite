package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
	serverstore "github.com/tronicum/punchbag-cube-testsuite/server/store"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockStore implements the Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateCluster(cluster *sharedmodels.Cluster) error {
	args := m.Called(cluster)
	return args.Error(0)
}

func (m *MockStore) GetCluster(id string) (*sharedmodels.Cluster, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedmodels.Cluster), args.Error(1)
}

func (m *MockStore) ListClusters() ([]*sharedmodels.Cluster, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.Cluster), args.Error(1)
}

func (m *MockStore) ListClustersByProvider(provider sharedmodels.CloudProvider) ([]*sharedmodels.Cluster, error) {
	args := m.Called(provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.Cluster), args.Error(1)
}

func (m *MockStore) UpdateCluster(id string, cluster *sharedmodels.Cluster) error {
	args := m.Called(id, cluster)
	return args.Error(0)
}

func (m *MockStore) DeleteCluster(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateTestResult(result *sharedmodels.TestResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockStore) GetTestResult(id string) (*sharedmodels.TestResult, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedmodels.TestResult), args.Error(1)
}

func (m *MockStore) ListTestResults(clusterID string) ([]*sharedmodels.TestResult, error) {
	args := m.Called(clusterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.TestResult), args.Error(1)
}

func (m *MockStore) UpdateTestResult(result *sharedmodels.TestResult) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *MockStore) CreateNodePool(nodePool *sharedmodels.NodePool) error {
	args := m.Called(nodePool)
	return args.Error(0)
}

func (m *MockStore) GetNodePool(id string) (*sharedmodels.NodePool, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedmodels.NodePool), args.Error(1)
}

func (m *MockStore) ListNodePools(clusterID string) ([]*sharedmodels.NodePool, error) {
	args := m.Called(clusterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.NodePool), args.Error(1)
}

func (m *MockStore) UpdateNodePool(id string, nodePool *sharedmodels.NodePool) error {
	args := m.Called(id, nodePool)
	return args.Error(0)
}

func (m *MockStore) DeleteNodePool(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureBudget(budget *sharedmodels.AzureBudget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockStore) GetAzureBudget(id string) (*sharedmodels.AzureBudget, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedmodels.AzureBudget), args.Error(1)
}

func (m *MockStore) ListAzureBudgets() ([]*sharedmodels.AzureBudget, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.AzureBudget), args.Error(1)
}

func (m *MockStore) UpdateAzureBudget(id string, budget *sharedmodels.AzureBudget) error {
	args := m.Called(id, budget)
	return args.Error(0)
}

func (m *MockStore) DeleteAzureBudget(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureMonitoring(monitoring *sharedmodels.AzureMonitoring) error {
	args := m.Called(monitoring)
	return args.Error(0)
}

func (m *MockStore) GetAzureMonitoring(id string) (*sharedmodels.AzureMonitoring, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedmodels.AzureMonitoring), args.Error(1)
}

func (m *MockStore) ListAzureMonitorings() ([]*sharedmodels.AzureMonitoring, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.AzureMonitoring), args.Error(1)
}

func (m *MockStore) UpdateAzureMonitoring(id string, monitoring *sharedmodels.AzureMonitoring) error {
	args := m.Called(id, monitoring)
	return args.Error(0)
}

func (m *MockStore) DeleteAzureMonitoring(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) CreateAzureKubernetes(kubernetes *sharedmodels.AzureKubernetes) error {
	args := m.Called(kubernetes)
	return args.Error(0)
}

func (m *MockStore) GetAzureKubernetes(id string) (*sharedmodels.AzureKubernetes, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedmodels.AzureKubernetes), args.Error(1)
}

func (m *MockStore) ListAzureKubernetes() ([]*sharedmodels.AzureKubernetes, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedmodels.AzureKubernetes), args.Error(1)
}

func (m *MockStore) UpdateAzureKubernetes(id string, kubernetes *sharedmodels.AzureKubernetes) error {
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
		cluster := sharedmodels.Cluster{
			Name:          "test-cluster",
			Provider: models.Azure,
			Status:        models.ClusterStatusCreating,
		}

		mockStore.On("CreateCluster", mock.AnythingOfType("*sharedmodels.Cluster")).Return(nil)

		body, _ := json.Marshal(cluster)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/clusters", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response sharedmodels.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, cluster.Name, response.Name)
		assert.NotEmpty(t, response.ID) // ID should be generated
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster already exists", func(t *testing.T) {
		cluster := sharedmodels.Cluster{
			Name:          "existing-cluster",
			Provider: models.Azure,
			Status:        models.ClusterStatusCreating,
		}

		mockStore.On("CreateCluster", mock.AnythingOfType("*sharedmodels.Cluster")).Return(serverstore.ErrAlreadyExists)

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
		cluster := &sharedmodels.Cluster{
			ID:            "cluster-1",
			Name:          "test-cluster",
			Provider: models.Azure,
			Status:        models.ClusterStatusRunning,
			CreatedAt:     time.Now(),
		}

		mockStore.On("GetCluster", "cluster-1").Return(cluster, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters/cluster-1", nil)
		
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response sharedmodels.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, cluster.ID, response.ID)
		assert.Equal(t, cluster.Name, response.Name)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		mockStore.On("GetCluster", "nonexistent").Return(nil, serverstore.ErrNotFound)

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
		clusters := []*sharedmodels.Cluster{
			{
				ID:            "cluster-1",
				Name:          "test-cluster-1",
				Provider: models.Azure,
				Status:        models.ClusterStatusRunning,
			},
			{
				ID:            "cluster-2",
				Name:          "test-cluster-2",
				Provider: models.AWS,
				Status:        models.ClusterStatusRunning,
			},
		}

		mockStore.On("ListClusters").Return(clusters, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters", nil)
		
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []*sharedmodels.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		mockStore.AssertExpectations(t)
	})

	t.Run("list clusters by provider", func(t *testing.T) {
		azureClusters := []*sharedmodels.Cluster{
			{
				ID:            "cluster-1",
				Name:          "azure-cluster",
				Provider: models.Azure,
				Status:        models.ClusterStatusRunning,
			},
		}

		mockStore.On("ListClustersByProvider", sharedmodels.CloudProvider("azure")).Return(azureClusters, nil)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/clusters?provider=azure", nil)
		
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []*sharedmodels.Cluster
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
		existingCluster := &sharedmodels.Cluster{
			ID:            "cluster-1",
			Name:          "old-name",
			Provider: models.Azure,
			Status:        models.ClusterStatusRunning,
		}

		updateData := sharedmodels.Cluster{
			Name:   "new-name",
			Status: "updating",
		}

		mockStore.On("GetCluster", "cluster-1").Return(existingCluster, nil)
		mockStore.On("UpdateCluster", mock.AnythingOfType("string"), mock.AnythingOfType("*sharedmodels.Cluster")).Return(nil)

		body, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("PUT", "/clusters/cluster-1", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response sharedmodels.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "new-name", response.Name)
		assert.Equal(t, "updating", response.Status)
		mockStore.AssertExpectations(t)
	})

	t.Run("cluster not found", func(t *testing.T) {
		updateData := sharedmodels.Cluster{Name: "new-name"}

		mockStore.On("GetCluster", "nonexistent").Return(nil, serverstore.ErrNotFound)

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
		mockStore.On("DeleteCluster", "nonexistent").Return(serverstore.ErrNotFound)

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
		cluster := &sharedmodels.Cluster{
			ID:            "cluster-1",
			Name:          "test-cluster",
			Provider: models.Azure,
			Status:        models.ClusterStatusRunning,
		}

		testRequest := models.TestRequest{
			ClusterID: "cluster-1",
			TestType:  "connectivity",
			Config:    map[string]interface{}{"timeout": "30s"},
		}

		mockStore.On("GetCluster", "cluster-1").Return(cluster, nil)
		mockStore.On("CreateTestResult", mock.AnythingOfType("*sharedmodels.TestResult")).Return(nil)

		body, _ := json.Marshal(testRequest)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)
		
		var response sharedmodels.TestResult
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

		mockStore.On("GetCluster", "nonexistent").Return(nil, serverstore.ErrNotFound)

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
		testResult := &sharedmodels.TestResult{
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
		
		var response sharedmodels.TestResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testResult.ID, response.ID)
		assert.Equal(t, testResult.Status, response.Status)
		mockStore.AssertExpectations(t)
	})

	t.Run("test result not found", func(t *testing.T) {
		mockStore.On("GetTestResult", "nonexistent").Return(nil, serverstore.ErrNotFound)

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
		testResults := []*sharedmodels.TestResult{
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
		
		var response []*sharedmodels.TestResult
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
		
		cluster := &sharedmodels.Cluster{
			ID:       "existing-cluster",
			Name:     "Test Cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusCreating,
		}

		mockStore.On("CreateCluster", mock.AnythingOfType("*sharedmodels.Cluster")).Return(serverstore.ErrAlreadyExists)

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
		
		azureClusters := []*sharedmodels.Cluster{
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
		
		var response []*sharedmodels.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, models.Azure, response[0].Provider)
		
		mockStore.AssertExpectations(t)
	})

	t.Run("run test on non-running cluster", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.POST("/clusters/:id/tests", handlers.RunTest)
		
		cluster := &sharedmodels.Cluster{
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
		mockStore.On("CreateCluster", mock.AnythingOfType("*sharedmodels.Cluster")).Return(nil).Times(5)

		// Create 5 concurrent requests
		results := make(chan int, 5)
		
		for i := 0; i < 5; i++ {
			go func(index int) {
				cluster := &sharedmodels.Cluster{
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

		cluster := &sharedmodels.Cluster{
			ID:       "test-cluster",
			Name:     "Test Cluster",
			Provider: models.Azure,
			Status:   models.ClusterStatusRunning,
		}

		// Mock cluster exists and test results can be created
		mockStore.On("GetCluster", "test-cluster").Return(cluster, nil).Times(3)
		mockStore.On("CreateTestResult", mock.AnythingOfType("*sharedmodels.TestResult")).Return(nil).Times(3)

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
		clusters := make([]*sharedmodels.Cluster, 1000)
		for i := 0; i < 1000; i++ {
			clusters[i] = &sharedmodels.Cluster{
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
		
		var response []*sharedmodels.Cluster
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1000)
		
		mockStore.AssertExpectations(t)
	})

	t.Run("large test result list", func(t *testing.T) {
		router, mockStore, handlers := setupTestRouter()
		router.GET("/clusters/:id/tests", handlers.ListTestResults)

		// Create a large list of test results
		testResults := make([]*sharedmodels.TestResult, 500)
		for i := 0; i < 500; i++ {
			testResults[i] = &sharedmodels.TestResult{
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
		
		var response []*sharedmodels.TestResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 500)
		
		mockStore.AssertExpectations(t)
	})
}

func TestProxyS3(t *testing.T) {
	router, _, handlers := setupTestRouter()
	router.POST("/api/proxy/aws/s3", handlers.ProxyS3)

	t.Run("valid creation with advanced fields", func(t *testing.T) {
		bucket := map[string]interface{}{
			"name": "test-bucket",
			"region": "us-east-1",
			"provider": "aws",
			"policy": map[string]interface{}{
				"version": "2012-10-17",
				"statement": []map[string]interface{}{
					{"effect": "Allow", "action": []string{"s3:GetObject"}, "resource": []string{"*"}, "principal": map[string]interface{}{"AWS": "*"}},
				},
			},
			"versioning": map[string]interface{}{"enabled": true},
			"lifecycle": []map[string]interface{}{
				{"id": "rule1", "status": "Enabled", "expiration_days": 30},
			},
		}
		body, _ := json.Marshal(bucket)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/proxy/aws/s3", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "test-bucket-aws", resp["id"])
		assert.Equal(t, "test-bucket", resp["name"])
		assert.Equal(t, "aws", resp["provider"])
		assert.NotEmpty(t, resp["created_at"])
	})

	t.Run("missing required fields", func(t *testing.T) {
		bucket := map[string]interface{}{"region": "us-east-1", "provider": "aws"}
		body, _ := json.Marshal(bucket)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/proxy/aws/s3", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "name, region, and provider are required")
	})

	t.Run("invalid policy", func(t *testing.T) {
		bucket := map[string]interface{}{
			"name": "test-bucket",
			"region": "us-east-1",
			"provider": "aws",
			"policy": map[string]interface{}{"version": ""},
		}
		body, _ := json.Marshal(bucket)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/proxy/aws/s3", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "policy.version and at least one statement required")
	})

	t.Run("invalid lifecycle rule", func(t *testing.T) {
		bucket := map[string]interface{}{
			"name": "test-bucket",
			"region": "us-east-1",
			"provider": "aws",
			"lifecycle": []map[string]interface{}{{"status": "Enabled"}},
		}
		body, _ := json.Marshal(bucket)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/proxy/aws/s3", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "lifecycle rule id and status required")
	})
}

func TestProxyObjectStorage_CreateBucket(t *testing.T) {
	router, _, handlers := setupTestRouter()
	router.POST("/api/proxy/:provider/objectstorage", handlers.ProxyObjectStorage)

	bucket := sharedmodels.ObjectStorageBucket{
		Name:     "test-bucket",
		Provider: sharedmodels.CloudProviderAWS,
		Region:   "eu-central-1",
		Policy: &sharedmodels.ObjectStoragePolicy{
			Version:   "2025-07-06",
			Statement: []sharedmodels.ObjectStorageStatement{{
				Effect:    "Allow",
				Principal: map[string]interface{}{"AWS": "*"},
				Action:    []string{"objectstorage:GetObject"},
				Resource:  []string{"arn:aws:s3:::test-bucket/*"},
			}},
		},
		Lifecycle: []sharedmodels.ObjectStorageRule{{
			ID:     "rule1",
			Status: "Enabled",
		}},
	}
	body, _ := json.Marshal(bucket)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/proxy/aws/objectstorage", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp sharedmodels.ObjectStorageBucket
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, bucket.Name, resp.Name)
	assert.Equal(t, bucket.Region, resp.Region)
	assert.Equal(t, bucket.Provider, resp.Provider)
	assert.NotEmpty(t, resp.ID)
	assert.NotZero(t, resp.CreatedAt)
	assert.NotNil(t, resp.Policy)
	assert.NotEmpty(t, resp.Lifecycle)
}
