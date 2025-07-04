package api

import (
	"net/http"
	"strconv"
	"time"

	"punchbag-cube-testsuite/server/models"
	"punchbag-cube-testsuite/server/store"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers contains the HTTP handlers for the API
type Handlers struct {
	store  store.Store
	logger *zap.Logger
}

// NewHandlers creates a new Handlers instance
func NewHandlers(store store.Store, logger *zap.Logger) *Handlers {
	return &Handlers{
		store:  store,
		logger: logger,
	}
}

// CreateCluster handles POST /clusters
func (h *Handlers) CreateCluster(c *gin.Context) {
	var cluster models.AKSCluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		h.logger.Error("Failed to bind cluster data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if not provided
	if cluster.ID == "" {
		cluster.ID = generateID()
	}

	if err := h.store.CreateCluster(&cluster); err != nil {
		if err == store.ErrAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "cluster already exists"})
			return
		}
		h.logger.Error("Failed to create cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("Cluster created", zap.String("id", cluster.ID))
	c.JSON(http.StatusCreated, cluster)
}

// GetCluster handles GET /clusters/:id
func (h *Handlers) GetCluster(c *gin.Context) {
	id := c.Param("id")
	cluster, err := h.store.GetCluster(id)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
		h.logger.Error("Failed to get cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, cluster)
}

// ListClusters handles GET /clusters
func (h *Handlers) ListClusters(c *gin.Context) {
	clusters, err := h.store.ListClusters()
	if err != nil {
		h.logger.Error("Failed to list clusters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clusters": clusters})
}

// UpdateCluster handles PUT /clusters/:id
func (h *Handlers) UpdateCluster(c *gin.Context) {
	id := c.Param("id")
	var cluster models.AKSCluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		h.logger.Error("Failed to bind cluster data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.UpdateCluster(id, &cluster); err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
		h.logger.Error("Failed to update cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("Cluster updated", zap.String("id", id))
	c.JSON(http.StatusOK, cluster)
}

// DeleteCluster handles DELETE /clusters/:id
func (h *Handlers) DeleteCluster(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.DeleteCluster(id); err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
		h.logger.Error("Failed to delete cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("Cluster deleted", zap.String("id", id))
	c.JSON(http.StatusNoContent, nil)
}

// RunTest handles POST /clusters/:id/tests
func (h *Handlers) RunTest(c *gin.Context) {
	clusterID := c.Param("id")
	var testReq models.AKSTestRequest
	if err := c.ShouldBindJSON(&testReq); err != nil {
		h.logger.Error("Failed to bind test request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify cluster exists
	_, err := h.store.GetCluster(clusterID)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
		h.logger.Error("Failed to get cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Create test result
	testResult := &models.AKSTestResult{
		ID:        generateID(),
		ClusterID: clusterID,
		TestType:  testReq.TestType,
		Status:    "running",
		Details:   testReq.Config,
	}

	if err := h.store.CreateTestResult(testResult); err != nil {
		h.logger.Error("Failed to create test result", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// In a real implementation, you would start the test asynchronously
	// For this example, we'll simulate test completion
	go h.simulateTest(testResult)

	h.logger.Info("Test started", zap.String("test_id", testResult.ID), zap.String("cluster_id", clusterID))
	c.JSON(http.StatusAccepted, testResult)
}

// GetTestResult handles GET /tests/:id
func (h *Handlers) GetTestResult(c *gin.Context) {
	id := c.Param("id")
	result, err := h.store.GetTestResult(id)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "test result not found"})
			return
		}
		h.logger.Error("Failed to get test result", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListTestResults handles GET /clusters/:id/tests
func (h *Handlers) ListTestResults(c *gin.Context) {
	clusterID := c.Param("id")
	results, err := h.store.ListTestResults(clusterID)
	if err != nil {
		h.logger.Error("Failed to list test results", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"test_results": results})
}

// simulateTest simulates a test execution and updates the result
func (h *Handlers) simulateTest(testResult *models.AKSTestResult) {
	// Simulate test duration
	time.Sleep(5 * time.Second)

	// Update test result
	now := time.Now()
	testResult.Status = "completed"
	testResult.Duration = now.Sub(testResult.StartedAt)
	testResult.CompletedAt = &now
	testResult.Details = map[string]interface{}{
		"requests_sent":       1000,
		"successful_requests": 995,
		"failed_requests":     5,
		"average_latency_ms":  45.2,
		"p95_latency_ms":      89.7,
		"p99_latency_ms":      156.3,
	}

	if err := h.store.UpdateTestResult(testResult.ID, testResult); err != nil {
		h.logger.Error("Failed to update test result", zap.Error(err))
	} else {
		h.logger.Info("Test completed", zap.String("test_id", testResult.ID))
	}
}

// generateID generates a simple ID for demonstration purposes
func generateID() string {
	return "id-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}
