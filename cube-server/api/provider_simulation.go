package api

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
	simulation "github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
	store "github.com/tronicum/punchbag-cube-testsuite/store"
	"go.uber.org/zap"
	"net/http"
)

// GenericAWSS3SimHandler simulates AWS S3 SDK-compatible endpoints for the provider 'generic-aws-s3'
func (h *ProviderSimulationHandlers) GenericAWSS3SimHandler(c *gin.Context) {
	// This is a stub for now. Log the request and return a generic simulated response.
	path := c.Param("path")
	method := c.Request.Method
	h.logger.Info("[GENERIC AWS S3 SIM] Simulating AWS S3 SDK endpoint", zap.String("method", method), zap.String("path", path))
	c.JSON(200, gin.H{"message": "Simulated AWS S3 SDK endpoint", "method": method, "path": path})
}

// ...existing code...

// ProviderSimulationHandlers contains simulation endpoints for cloud providers
type ProviderSimulationHandlers struct {
	store     store.Store
	logger    *zap.Logger
	simulator *simulation.SimulationService
}

// DeleteSimulatedBucket handles DELETE /api/v1/simulate/providers/:provider/buckets/:bucket
func (h *ProviderSimulationHandlers) DeleteSimulatedBucket(c *gin.Context) {
	provider := c.Param("provider")
	bucket := c.Param("bucket")

	simReq := &simulation.SimulationRequest{
		Provider:   provider,
		Operation:  "delete_bucket",
		Parameters: map[string]interface{}{"bucket": bucket},
	}
	result := h.simulator.SimulateOperation(simReq)
	if result.Success {
		c.JSON(http.StatusOK, result.Result)
	} else {
		c.JSON(http.StatusNotFound, result)
	}
}

// NewProviderSimulationHandlers creates a new ProviderSimulationHandlers instance
func NewProviderSimulationHandlers(s store.Store, logger *zap.Logger, sim *simulation.SimulationService) *ProviderSimulationHandlers {
	return &ProviderSimulationHandlers{
		store:     s,
		logger:    logger,
		simulator: sim,
	}
}

// CreateSimulatedBucket handles POST /api/v1/simulate/providers/:provider/buckets
func (h *ProviderSimulationHandlers) CreateSimulatedBucket(c *gin.Context) {
	// DEBUG: Log entry to handler
	fmt.Println("[DEBUG] Entered CreateSimulatedBucket handler")
	provider := c.Param("provider")
	var req map[string]interface{}
	fmt.Printf("[SERVER DEBUG] CreateSimulatedBucket called for provider=%s\n", provider)
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[SERVER DEBUG] Invalid request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Printf("[SERVER DEBUG] Request body: %#v\n", req)
	simReq := &simulation.SimulationRequest{
		Provider:   provider,
		Operation:  "create_bucket",
		Parameters: req,
	}
	fmt.Printf("[SERVER DEBUG] Calling SimulateOperation with: provider=%s, op=%s, params=%#v\n", simReq.Provider, simReq.Operation, simReq.Parameters)
	result := h.simulator.SimulateOperation(simReq)
	fmt.Printf("[SERVER DEBUG] SimulateOperation result: success=%v, result=%#v, error=%v\n", result.Success, result.Result, result.Error)
	if result.Success {
		c.JSON(http.StatusCreated, result.Result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// ListSimulatedBuckets handles GET /api/v1/simulate/providers/:provider/buckets
func (h *ProviderSimulationHandlers) ListSimulatedBuckets(c *gin.Context) {
	fmt.Println("[DEBUG] Entered ListSimulatedBuckets handler")
	provider := c.Param("provider")
	fmt.Printf("[SERVER DEBUG] ListSimulatedBuckets called for provider=%s\n", provider)
	simReq := &simulation.SimulationRequest{
		Provider:   provider,
		Operation:  "list_buckets",
		Parameters: map[string]interface{}{},
	}
	fmt.Printf("[SERVER DEBUG] Calling SimulateOperation for list_buckets with: provider=%s\n", provider)
	result := h.simulator.SimulateOperation(simReq)
	fmt.Printf("[SERVER DEBUG] SimulateOperation list_buckets result: success=%v, result=%#v, error=%v\n", result.Success, result.Result, result.Error)
	if result.Success {
		// result.Result is a map with "buckets": []map[string]interface{}
		if bucketsObj, ok := result.Result["buckets"]; ok {
			if bucketsList, ok := bucketsObj.([]map[string]interface{}); ok {
				// Convert simulation bucket format to ObjectStorageBucket format
				var convertedBuckets []sharedmodels.ObjectStorageBucket
				for _, simBucket := range bucketsList {
					bucket := sharedmodels.ObjectStorageBucket{
						Name:     getString(simBucket, "bucket"),
						Provider: sharedmodels.CloudProvider(getString(simBucket, "provider")),
						Region:   getString(simBucket, "region"),
					}
					// Set a default CreatedAt if not present
					if bucket.Name != "" {
						bucket.CreatedAt = time.Now().UTC()
					}
					convertedBuckets = append(convertedBuckets, bucket)
				}
				fmt.Printf("[SERVER DEBUG] Converted %d buckets to ObjectStorageBucket format\n", len(convertedBuckets))
				c.JSON(http.StatusOK, convertedBuckets)
			} else {
				fmt.Printf("[SERVER DEBUG] bucketsObj is not []map[string]interface{}, type=%T\n", bucketsObj)
				c.JSON(http.StatusOK, []interface{}{})
			}
		} else {
			fmt.Printf("[SERVER DEBUG] No 'buckets' key in result\n")
			c.JSON(http.StatusOK, []interface{}{})
		}
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// ValidateProvider handles POST /api/v1/providers/validate
func (h *ProviderSimulationHandlers) ValidateProvider(c *gin.Context) {
	var req struct {
		Provider    string                 `json:"provider" binding:"required"`
		Credentials map[string]interface{} `json:"credentials,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	h.logger.Info("Validating provider", zap.String("provider", req.Provider))

	result := h.simulator.ValidateProvider(req.Provider, req.Credentials)

	if result.Valid {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// ValidateProviderLegacy handles GET /validate/{provider-name}/ (legacy endpoint)
// Deprecated: Use ValidateProvider instead. This will be removed in a future release.
func (h *ProviderSimulationHandlers) ValidateProviderLegacy(c *gin.Context) {
	provider := c.Param("provider")

	h.logger.Info("Validating provider (legacy)", zap.String("provider", provider))

	result := h.simulator.ValidateProvider(provider, nil)

	if result.Valid {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// The following provider-specific validation endpoints are deprecated and replaced by the shared simulation logic.
// They are removed to ensure all validation uses the shared simulation service.

// GetProviderInfo handles GET /providers/{provider-name}/info
func (h *ProviderSimulationHandlers) GetProviderInfo(c *gin.Context) {
	provider := c.Param("provider")

	switch provider {
	case "azure":
		h.getAzureInfo(c)
	case "hetzner-hcloud":
		h.getHetznerInfo(c)
	case "united-ionos":
		h.getIONOSInfo(c)
	case "schwarz-stackit":
		h.getStackITInfo(c)
	case "aws":
		h.getAWSInfo(c)
	case "gcp":
		h.getGCPInfo(c)
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("provider not found: %s", provider),
		})
	}
}

// The following provider-specific validation endpoints are deprecated and replaced by the shared simulation logic.
// They are removed to ensure all validation uses the shared simulation service.
//
// func (h *ProviderSimulationHandlers) validateAzure(c *gin.Context) { ... }
// func (h *ProviderSimulationHandlers) validateHetzner(c *gin.Context) { ... }
// func (h *ProviderSimulationHandlers) validateIONOS(c *gin.Context) { ... }
// func (h *ProviderSimulationHandlers) validateStackIT(c *gin.Context) { ... }
// func (h *ProviderSimulationHandlers) validateAWS(c *gin.Context) { ... }
// func (h *ProviderSimulationHandlers) validateGCP(c *gin.Context) { ... }
//
// Please use ValidateProvider instead.

func (h *ProviderSimulationHandlers) getAzureInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provider":      "azure",
		"name":          "Microsoft Azure",
		"description":   "Microsoft's cloud computing platform",
		"documentation": "https://docs.microsoft.com/en-us/azure/aks/",
		"pricing_model": "pay-as-you-go",
		"supported_features": []string{
			"auto-scaling",
			"load-balancing",
			"monitoring",
			"rbac",
			"network-policies",
		},
	})
}

func (h *ProviderSimulationHandlers) getHetznerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provider":      "hetzner-hcloud",
		"name":          "Hetzner Cloud",
		"description":   "German cloud hosting provider with competitive pricing",
		"documentation": "https://docs.hetzner.com/cloud/",
		"pricing_model": "hourly-billing",
		"supported_features": []string{
			"auto-scaling",
			"load-balancing",
			"private-networks",
			"ssh-keys",
			"firewalls",
		},
	})
}

func (h *ProviderSimulationHandlers) getIONOSInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provider":      "united-ionos",
		"name":          "IONOS Cloud",
		"description":   "European cloud provider with data sovereignty focus",
		"documentation": "https://docs.ionos.com/cloud/",
		"pricing_model": "hourly-billing",
		"supported_features": []string{
			"kubernetes",
			"managed-services",
			"data-sovereignty",
			"compliance",
			"monitoring",
		},
	})
}

func (h *ProviderSimulationHandlers) getStackITInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provider":      "schwarz-stackit",
		"name":          "StackIT",
		"description":   "Schwarz Group's cloud platform for enterprise customers",
		"documentation": "https://docs.stackit.cloud/",
		"pricing_model": "enterprise-contracts",
		"supported_features": []string{
			"ske",
			"enterprise-grade",
			"compliance",
			"private-cloud",
			"managed-kubernetes",
		},
	})
}

func (h *ProviderSimulationHandlers) getAWSInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provider":      "aws",
		"name":          "Amazon Web Services",
		"description":   "Amazon's comprehensive cloud computing platform",
		"documentation": "https://docs.aws.amazon.com/eks/",
		"pricing_model": "pay-as-you-go",
		"supported_features": []string{
			"eks",
			"fargate",
			"auto-scaling",
			"iam-integration",
			"vpc-native",
		},
	})
}

func (h *ProviderSimulationHandlers) getGCPInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provider":      "gcp",
		"name":          "Google Cloud Platform",
		"description":   "Google's cloud computing services",
		"documentation": "https://cloud.google.com/kubernetes-engine/docs",
		"pricing_model": "pay-as-you-go",
		"supported_features": []string{
			"gke",
			"autopilot",
			"workload-identity",
			"istio-integration",
			"binary-authorization",
		},
	})
}

// SimulateProviderOperation handles POST /api/v1/providers/simulate
// Uses the shared simulation service for all provider operations.
func (h *ProviderSimulationHandlers) SimulateProviderOperation(c *gin.Context) {
	var req simulation.SimulationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	h.logger.Info("Simulating provider operation",
		zap.String("provider", req.Provider),
		zap.String("operation", req.Operation))

	result := h.simulator.SimulateOperation(&req)

	if result.Success {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, result)
	}
}

// CreateSimulatedCluster creates a simulated cluster using the shared simulation service
func (h *ProviderSimulationHandlers) CreateSimulatedCluster(c *gin.Context) {
	var req sharedmodels.ClusterCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	h.logger.Info("Creating simulated cluster",
		zap.String("name", req.Name),
		zap.String("provider", string(req.Provider)))

	// Generate simulated cluster using shared service
	cluster := h.simulator.GenerateClusterFromSimulation(string(req.Provider), req.Name, req.Config)

	// Store the simulated cluster
	_, err := h.store.CreateCluster(cluster)
	if err != nil {
		h.logger.Error("Failed to store simulated cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create simulated cluster",
		})
		return
	}

	c.JSON(http.StatusCreated, cluster)
}

// RunSimulatedTest runs a simulated test using the shared simulation service
func (h *ProviderSimulationHandlers) RunSimulatedTest(c *gin.Context) {
	var req sharedmodels.TestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	h.logger.Info("Running simulated test",
		zap.String("cluster_id", req.ClusterID),
		zap.String("test_type", req.TestType))

	// Check if cluster exists
	_, err := h.store.GetCluster(req.ClusterID)
	if err != nil {
		if err != nil && (err.Error() == "cluster not found" || err.Error() == "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cluster not found",
			})
			return
		}
		h.logger.Error("Failed to get cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get cluster",
		})
		return
	}

	// Generate simulated test result using shared service
	testResult := h.simulator.GenerateTestResultFromSimulation(req.ClusterID, req.TestType)

	// Store the test result
	_, err = h.store.CreateTestResult(testResult)
	if err != nil {
		h.logger.Error("Failed to store test result", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create test result",
		})
		return
	}

	c.JSON(http.StatusCreated, testResult)
}

// API DOCS: Multicloud Simulation
//
// All simulation endpoints accept a 'provider' parameter (e.g., 'azure', 'aws', 'gcp').
// Example: POST /api/v1/providers/:provider/operations/:operation
//
// Request body for cluster creation:
// {
//   "provider": "aws",
//   "operation": "create_cluster",
//   "parameters": {
//     "name": "my-eks",
//     "region": "us-west-2",
//     "node_count": 3,
//     ...
//   }
// }
//
// The response will include provider-specific fields. See simulation/service.go for details.
//
// TODO: Document supported resource types and required fields for each provider.

// getString safely extracts a string value from a map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}