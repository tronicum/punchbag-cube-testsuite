package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tronicum/punchbag-cube-testsuite/store"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"

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

	// Generate ID if not provided
	if cluster.ID == "" {
		cluster.ID = generateID()
	}

	// Validate provider
	if cluster.Provider == "" {
		// Fix: use sharedmodels.CloudProvider instead of sharedmodels.Provider
		cluster.Provider = sharedmodels.CloudProvider("azure") // default to Azure for backward compatibility
	}

	// Validate required fields based on provider
	if err := h.validateClusterByProvider(&cluster); err != nil {
		h.logger.Error("Cluster validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.store.CreateCluster(&cluster)
	if err != nil {
		if err != nil && (err.Error() == "cluster already exists" || err.Error() == "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "cluster already exists"})
			return
		}
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
		if err != nil && (err.Error() == "cluster not found" || err.Error() == "not found") {
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
		// Filter by provider
		clustProv := sharedmodels.CloudProvider(provider)
		clusters, err = h.store.ListClustersByProvider(clustProv)
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
	var cluster sharedmodels.Cluster
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

	updated, err := h.store.UpdateCluster(id, &cluster)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
			return
		}
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
	err := h.store.DeleteCluster(id)
	if err != nil {
		if err != nil && (err.Error() == "cluster not found" || err.Error() == "not found") {
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
	cluster, err := h.store.GetCluster(clusterID)
	if err != nil {
		if err != nil && (err.Error() == "cluster not found" || err.Error() == "not found") {
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

	// Add provider information to test details
	if testResult.Details == nil {
		testResult.Details = make(map[string]interface{})
	}
	testResult.Details["provider"] = string(cluster.Provider)
	testResult.Details["cluster_name"] = cluster.Name

	createdTest, err := h.store.CreateTestResult(testResult)
	if err != nil {
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
	c.JSON(http.StatusAccepted, createdTest)
}

// GetTestResult handles GET /tests/:id
func (h *Handlers) GetTestResult(c *gin.Context) {
	id := c.Param("id")
	result, err := h.store.GetTestResult(id)
	if err != nil {
		if err != nil && (err.Error() == "test result not found" || err.Error() == "not found") {
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

// ProxyS3 handles /api/proxy/:provider/s3
func (h *Handlers) ProxyS3(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var bucket sharedmodels.ObjectStorageBucket
		if err := c.ShouldBindJSON(&bucket); err != nil {
			h.logger.Error("Failed to bind S3 bucket payload", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if bucket.Name == "" || bucket.Region == "" || bucket.Provider == "" {
			h.logger.Error("Missing required S3 bucket fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, region, and provider are required"})
			return
		}
		bucket.ID = bucket.Name + "-" + string(bucket.Provider)
		bucket.CreatedAt = time.Now()
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

// ProxyObjectStorage handles /api/proxy/:provider/objectstorage for all providers
func (h *Handlers) ProxyObjectStorage(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var bucket sharedmodels.ObjectStorageBucket
		if err := c.ShouldBindJSON(&bucket); err != nil {
			h.logger.Error("Failed to bind object storage bucket payload", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if bucket.Name == "" || bucket.Provider == "" {
			h.logger.Error("Missing required object storage bucket fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "name and provider are required"})
			return
		}
		bucket.ID = bucket.Name + "-" + string(bucket.Provider)
		bucket.CreatedAt = time.Now()
		bucket.UpdatedAt = time.Now()
		// Provider-specific logic can be added here if needed
		switch bucket.Provider {
		case "stackit":
			bucket.ProviderConfig["note"] = "StackIT object storage bucket created (mock)"
		case "hetzner":
			bucket.ProviderConfig["note"] = "Hetzner object storage bucket created (mock)"
		case "ionos":
			bucket.ProviderConfig["note"] = "IONOS object storage bucket created (mock)"
		}
		c.JSON(http.StatusCreated, bucket)
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Object storage proxy only supports POST for now"})
}

// GetCloudFormationStack handles GET /api/v1/cloudformation/stack?name=<name>
func (h *Handlers) GetCloudFormationStack(c *gin.Context) {
	stackName := c.Query("name")
	if stackName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stack name is required"})
		return
	}
	// Simulate: return a minimal template (in real use, fetch from AWS or mock DB)
	template := `AWSTemplateFormatVersion: '2010-09-09'
Description: Simulated CloudFormation template
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties: {}`
	c.Data(http.StatusOK, "application/x-yaml", []byte(template))
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

	// Add provider-specific test results
	provider := testResult.Details["provider"].(string)
	testResult.Details["requests_sent"] = 1000
	testResult.Details["successful_requests"] = 995
	testResult.Details["failed_requests"] = 5
	testResult.Details["average_latency_ms"] = 45.2
	testResult.Details["p95_latency_ms"] = 89.7
	testResult.Details["p99_latency_ms"] = 156.3
	testResult.Details["provider_specific"] = h.getProviderSpecificMetrics(provider)

	_, err := h.store.UpdateTestResult(testResult.ID, testResult)
	if err != nil {
		h.logger.Error("Failed to update test result", zap.Error(err))
	} else {
		h.logger.Info("Test completed", zap.String("test_id", testResult.ID))
	}
}

// getProviderSpecificMetrics returns provider-specific test metrics
func (h *Handlers) getProviderSpecificMetrics(provider string) map[string]interface{} {
	switch provider {
	case string(sharedmodels.CloudProviderStackIT):
		return map[string]interface{}{
			"stackit_specific_metric": "sample_value",
			"project_usage": "normal",
			"hibernation_supported": true,
		}
	case string(sharedmodels.CloudProviderHetzner):
		return map[string]interface{}{
			"hetzner_specific_metric": "sample_value",
			"network_zone": "eu-central",
			"load_balancer_type": "lb11",
			"server_type": "cx21",
		}
	case string(sharedmodels.CloudProviderIONOS):
		return map[string]interface{}{
			"ionos_specific_metric": "sample_value",
			"datacenter_location": "de/fra",
			"k8s_cluster_type": "managed",
			"maintenance_window": "automatic",
		}
	case string(sharedmodels.CloudProviderAzure):
		return map[string]interface{}{
			"azure_specific_metric": "sample_value",
			"aks_version": "1.28.0",
			"resource_group_location": "eastus",
		}
	case string(sharedmodels.CloudProviderAWS):
		return map[string]interface{}{
			"aws_specific_metric": "sample_value",
			"eks_version": "1.28",
			"vpc_configuration": "standard",
		}
	case string(sharedmodels.CloudProviderGCP):
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
func (h *Handlers) validateClusterByProvider(cluster *sharedmodels.Cluster) error {
	switch cluster.Provider {
	case sharedmodels.CloudProviderStackIT:
		if cluster.ProjectID == "" {
			return fmt.Errorf("project_id is required for StackIT clusters")
		}
		if cluster.Location == "" && cluster.Region == "" {
			return fmt.Errorf("location or region is required for StackIT clusters")
		}
	case sharedmodels.CloudProviderHetzner:
		if cluster.Location == "" {
			return fmt.Errorf("location is required for Hetzner Cloud clusters")
		}
		// Validate Hetzner-specific fields from ProviderConfig
		if cluster.ProviderConfig != nil {
			if hetznerConfig, ok := cluster.ProviderConfig["hetzner_config"]; ok {
				if config, ok := hetznerConfig.(map[string]interface{}); ok {
					// Check server_type field
					if serverType, exists := config["server_type"]; exists && serverType == "" {
						return fmt.Errorf("server_type cannot be empty for Hetzner clusters")
					}
				}
			}
		}
	case sharedmodels.CloudProviderIONOS:
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
	case sharedmodels.CloudProviderAzure:
		if cluster.ResourceGroup == "" {
			return fmt.Errorf("resource_group is required for Azure clusters")
		}
		if cluster.Location == "" {
			return fmt.Errorf("location is required for Azure clusters")
		}
	case sharedmodels.CloudProviderAWS:
		if cluster.Region == "" {
			return fmt.Errorf("region is required for AWS clusters")
		}
	case sharedmodels.CloudProviderGCP:
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
