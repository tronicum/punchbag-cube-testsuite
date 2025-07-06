package api

import (
	"net/http"
	"strconv"
	"time"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
	"github.com/tronicum/punchbag-cube-testsuite/store"

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
	var cluster sharedmodels.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		h.logger.Error("Failed to bind cluster data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cloud provider
	validProviders := []sharedmodels.CloudProvider{
		sharedmodels.Azure,
		sharedmodels.AWS,
		sharedmodels.GCP,
		sharedmodels.StackIT,
		sharedmodels.Hetzner,
		sharedmodels.IONOS,
	}
	valid := false
	for _, provider := range validProviders {
		if cluster.Provider == provider {
			valid = true
			break
		}
	}
	if !valid {
		h.logger.Error("Invalid cloud provider", zap.String("provider", string(cluster.Provider)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cloud provider"})
		return
	}

	// Generate ID if not provided
	if cluster.ID == "" {
		cluster.ID = generateID()
	}

	// Convert to server model and create
	created, err := h.store.CreateCluster(&cluster)
	if err != nil {
		h.logger.Error("Failed to create cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("Cluster created", zap.String("id", created.ID), zap.String("provider", string(created.Provider)))
	c.JSON(http.StatusCreated, created)
}

// GetCluster handles GET /clusters/:id
func (h *Handlers) GetCluster(c *gin.Context) {
	id := c.Param("id")
	cluster, err := h.store.GetCluster(id)
	if err != nil {
		if err.Error() == "cluster not found" {
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
	provider := c.Query("provider")

	var clusters []*sharedmodels.Cluster
	var err error

	if provider != "" {
		clusters, err = h.store.ListClustersByProvider(sharedmodels.CloudProvider(provider))
	} else {
		clusters, err = h.store.ListClusters()
	}

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
	var cluster sharedmodels.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		h.logger.Error("Failed to bind cluster data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.store.UpdateCluster(id, &cluster)
	if err != nil {
		h.logger.Error("Failed to update cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("Cluster updated", zap.String("id", id), zap.String("provider", string(updated.Provider)))
	c.JSON(http.StatusOK, updated)
}

// DeleteCluster handles DELETE /clusters/:id
func (h *Handlers) DeleteCluster(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.DeleteCluster(id); err != nil {
		if err.Error() == "cluster not found" {
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
	var testReq sharedmodels.TestRequest
	if err := c.ShouldBindJSON(&testReq); err != nil {
		h.logger.Error("Failed to bind test request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify cluster exists
	_, err := h.store.GetCluster(clusterID)
	if err != nil {
		if err.Error() == "cluster not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
		h.logger.Error("Failed to get cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Create test result
	testResult := &sharedmodels.TestResult{
		ID:        generateID(),
		ClusterID: clusterID,
		TestType:  testReq.TestType,
		Status:    "running",
		Details:   testReq.Config,
	}

	created, err := h.store.CreateTestResult(testResult)
	if err != nil {
		h.logger.Error("Failed to create test result", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// In a real implementation, you would start the test asynchronously
	// For this example, we'll simulate test completion
	go h.simulateTest(created)

	h.logger.Info("Test started", zap.String("test_id", created.ID), zap.String("cluster_id", clusterID))
	c.JSON(http.StatusAccepted, created)
}

// GetTestResult handles GET /tests/:id
func (h *Handlers) GetTestResult(c *gin.Context) {
	id := c.Param("id")
	result, err := h.store.GetTestResult(id)
	if err != nil {
		if err.Error() == "test result not found" {
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

// ProxyS3 handles /api/proxy/aws/s3
func (h *Handlers) ProxyS3(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var bucket sharedmodels.ObjectStorageBucket
		if err := c.ShouldBindJSON(&bucket); err != nil {
			h.logger.Error("Failed to bind S3 bucket payload", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// --- Validation for advanced S3 features ---
		if bucket.Name == "" || bucket.Region == "" || bucket.Provider == "" {
			h.logger.Error("Missing required S3 bucket fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, region, and provider are required"})
			return
		}
		if bucket.Policy != nil {
			if bucket.Policy.Version == "" || len(bucket.Policy.Statement) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "policy.version and at least one statement required"})
				return
			}
		}
		if bucket.Lifecycle != nil {
			for _, rule := range bucket.Lifecycle {
				if rule.ID == "" || rule.Status == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "lifecycle rule id and status required"})
					return
				}
			}
		}

		// --- Simulate creation/response ---
		bucket.ID = bucket.Name + "-" + string(bucket.Provider)
		bucket.CreatedAt = time.Now()

		// TODO: Integrate with real provider APIs here

		c.JSON(http.StatusCreated, bucket)
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{"message": "S3 proxy only supports POST for now"})
}

// ProxyBlob handles /api/proxy/azure/blob
func (h *Handlers) ProxyBlob(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Blob proxy not implemented yet"})
}

// ProxyGCS handles /api/proxy/gcp/gcs
func (h *Handlers) ProxyGCS(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GCS proxy not implemented yet"})
}

// simulateTest simulates a test execution and updates the result
func (h *Handlers) simulateTest(testResult *sharedmodels.TestResult) {
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

	_, err := h.store.UpdateTestResult(testResult.ID, testResult)
	if err != nil {
		h.logger.Error("Failed to update test result", zap.Error(err))
	} else {
		h.logger.Info("Test completed", zap.String("test_id", testResult.ID))
	}
}

// generateID generates a simple ID for demonstration purposes
func generateID() string {
	return "id-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}
