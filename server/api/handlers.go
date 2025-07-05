package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/username/punchbag-cube-testsuite/server/models"
	"github.com/username/punchbag-cube-testsuite/server/store"

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
	var cluster models.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		h.logger.Error("Failed to bind cluster data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if not provided
	if cluster.ID == "" {
		cluster.ID = generateID()
	}

	// Validate provider
	if cluster.Provider == "" {
		cluster.Provider = models.CloudProviderAzure // default to Azure for backward compatibility
	}

	// Validate required fields based on provider
	if err := h.validateClusterByProvider(&cluster); err != nil {
		h.logger.Error("Cluster validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	h.logger.Info("Cluster created", zap.String("id", cluster.ID), zap.String("provider", string(cluster.Provider)))
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
	provider := c.Query("provider")
	
	var clusters []*models.Cluster
	var err error
	
	if provider != "" {
		// Filter by provider
		cloudProvider := models.CloudProvider(provider)
		clusters, err = h.store.ListClustersByProvider(cloudProvider)
	} else {
		// List all clusters
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
	var cluster models.Cluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		h.logger.Error("Failed to bind cluster data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate provider-specific fields
	if err := h.validateClusterByProvider(&cluster); err != nil {
		h.logger.Error("Cluster validation failed", zap.Error(err))
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

	h.logger.Info("Cluster updated", zap.String("id", id), zap.String("provider", string(cluster.Provider)))
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
	var testReq models.TestRequest
	if err := c.ShouldBindJSON(&testReq); err != nil {
		h.logger.Error("Failed to bind test request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify cluster exists
	cluster, err := h.store.GetCluster(clusterID)
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
	testResult := &models.TestResult{
		ID:        generateID(),
		ClusterID: clusterID,
		TestType:  testReq.TestType,
		Status:    "running",
		Details:   testReq.Config,
	}

	// Add provider information to test details
	if testResult.Details == nil {
		testResult.Details = make(map[string]interface{})
	}
	testResult.Details["provider"] = string(cluster.Provider)
	testResult.Details["cluster_name"] = cluster.Name

	if err := h.store.CreateTestResult(testResult); err != nil {
		h.logger.Error("Failed to create test result", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// In a real implementation, you would start the test asynchronously
	// For this example, we'll simulate test completion
	go h.simulateTest(testResult)

	h.logger.Info("Test started", 
		zap.String("test_id", testResult.ID), 
		zap.String("cluster_id", clusterID),
		zap.String("provider", string(cluster.Provider)))
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

// ValidateProvider handles GET /validate/:provider
func (h *Handlers) ValidateProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	// Simulate validation logic
	c.JSON(http.StatusOK, gin.H{
		"provider": provider,
		"status":   "valid",
	})
}

// SimulateProviderOperation handles POST /providers/:provider/operations/:operation
func (h *Handlers) SimulateProviderOperation(c *gin.Context) {
	provider := c.Param("provider")
	operation := c.Param("operation")
	if provider == "" || operation == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider and operation are required"})
		return
	}

	// Simulate operation logic
	c.JSON(http.StatusOK, gin.H{
		"provider":  provider,
		"operation": operation,
		"status":    "success",
	})
}

// simulateTest simulates a test execution and updates the result
func (h *Handlers) simulateTest(testResult *models.TestResult) {
	// Simulate test duration
	time.Sleep(5 * time.Second)

	// Update test result
	now := time.Now()
	testResult.Status = "completed"
	testResult.Duration = now.Sub(testResult.StartedAt)
	testResult.CompletedAt = &now
	
	// Add provider-specific test results
	provider := testResult.Details["provider"].(string)
	testResult.Details["requests_sent"] = 1000
	testResult.Details["successful_requests"] = 995
	testResult.Details["failed_requests"] = 5
	testResult.Details["average_latency_ms"] = 45.2
	testResult.Details["p95_latency_ms"] = 89.7
	testResult.Details["p99_latency_ms"] = 156.3
	testResult.Details["provider_specific"] = h.getProviderSpecificMetrics(provider)

	if err := h.store.UpdateTestResult(testResult.ID, testResult); err != nil {
		h.logger.Error("Failed to update test result", zap.Error(err))
	} else {
		h.logger.Info("Test completed", zap.String("test_id", testResult.ID))
	}
}

// getProviderSpecificMetrics returns provider-specific test metrics
func (h *Handlers) getProviderSpecificMetrics(provider string) map[string]interface{} {
	switch provider {
	case string(models.CloudProviderStackIT):
		return map[string]interface{}{
			"stackit_specific_metric": "sample_value",
			"project_usage": "normal",
			"hibernation_supported": true,
		}
	case string(models.CloudProviderHetzner):
		return map[string]interface{}{
			"hetzner_specific_metric": "sample_value",
			"network_zone": "eu-central",
			"load_balancer_type": "lb11",
			"server_type": "cx21",
		}
	case string(models.CloudProviderIONOS):
		return map[string]interface{}{
			"ionos_specific_metric": "sample_value",
			"datacenter_location": "de/fra",
			"k8s_cluster_type": "managed",
			"maintenance_window": "automatic",
		}
	case string(models.CloudProviderAzure):
		return map[string]interface{}{
			"azure_specific_metric": "sample_value",
			"aks_version": "1.28.0",
			"resource_group_location": "eastus",
		}
	case string(models.CloudProviderAWS):
		return map[string]interface{}{
			"aws_specific_metric": "sample_value",
			"eks_version": "1.28",
			"vpc_configuration": "standard",
		}
	case string(models.CloudProviderGCP):
		return map[string]interface{}{
			"gcp_specific_metric": "sample_value",
			"gke_version": "1.28.0",
			"autopilot_enabled": false,
		}
	default:
		return map[string]interface{}{}
	}
}

// validateClusterByProvider validates cluster configuration based on provider
func (h *Handlers) validateClusterByProvider(cluster *models.Cluster) error {
	switch cluster.Provider {
	case models.CloudProviderStackIT:
		if cluster.ProjectID == "" {
			return fmt.Errorf("project_id is required for StackIT clusters")
		}
		if cluster.Location == "" && cluster.Region == "" {
			return fmt.Errorf("location or region is required for StackIT clusters")
		}
	case models.CloudProviderHetzner:
		if cluster.Location == "" {
			return fmt.Errorf("location is required for Hetzner Cloud clusters")
		}
		// Validate Hetzner-specific fields from ProviderConfig
		if cluster.ProviderConfig != nil {
			if hetznerConfig, ok := cluster.ProviderConfig["hetzner_config"]; ok {
				if config, ok := hetznerConfig.(map[string]interface{}); ok {
					if serverType, exists := config["server_type"]; exists && serverType == "" {
						return fmt.Errorf("server_type cannot be empty for Hetzner clusters")
					}
				}
			}
		}
	case models.CloudProviderIONOS:
		// Validate IONOS-specific fields from ProviderConfig
		if cluster.ProviderConfig != nil {
			if ionosConfig, ok := cluster.ProviderConfig["ionos_config"]; ok {
				if config, ok := ionosConfig.(map[string]interface{}); ok {
					if datacenterID, exists := config["datacenter_id"]; !exists || datacenterID == "" {
						return fmt.Errorf("datacenter_id is required for IONOS Cloud clusters")
					}
				}
			} else {
				return fmt.Errorf("ionos_config is required for IONOS Cloud clusters")
			}
		} else {
			return fmt.Errorf("configuration is required for IONOS Cloud clusters")
		}
	case models.CloudProviderAzure:
		if cluster.ResourceGroup == "" {
			return fmt.Errorf("resource_group is required for Azure clusters")
		}
		if cluster.Location == "" {
			return fmt.Errorf("location is required for Azure clusters")
		}
	case models.CloudProviderAWS:
		if cluster.Region == "" {
			return fmt.Errorf("region is required for AWS clusters")
		}
	case models.CloudProviderGCP:
		if cluster.ProjectID == "" {
			return fmt.Errorf("project_id is required for GCP clusters")
		}
		if cluster.Region == "" {
			return fmt.Errorf("region is required for GCP clusters")
		}
	}
	return nil
}

// generateID generates a simple ID for demonstration purposes
func generateID() string {
	return "id-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}
