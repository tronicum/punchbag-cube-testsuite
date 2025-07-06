package api

import (
	"time"

	"github.com/username/punchbag-cube-testsuite/server/store"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes configures all the API routes
func SetupRoutes(router *gin.Engine, store store.Store, logger *zap.Logger) {
	handlers := NewHandlers(store, logger)

	// API version prefix
	v1 := router.Group("/api/v1")
	{
		// Cluster management endpoints
		clusters := v1.Group("/clusters")
		{
			clusters.POST("", handlers.CreateCluster)
			clusters.GET("", handlers.ListClusters)
			clusters.GET("/:id", handlers.GetCluster)
			clusters.PUT("/:id", handlers.UpdateCluster)
			clusters.DELETE("/:id", handlers.DeleteCluster)

			// Test endpoints for specific clusters
			clusters.POST("/:id/tests", handlers.RunTest)
			clusters.GET("/:id/tests", handlers.ListTestResults)
		}

		// Test result endpoints
		tests := v1.Group("/tests")
		{
			tests.GET("/:id", handlers.GetTestResult)
		}

		// Metrics and monitoring endpoints
		metrics := v1.Group("/metrics")
		{
			metrics.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"status":    "healthy",
					"timestamp": time.Now().Format(time.RFC3339),
				})
			})

			metrics.GET("/status", func(c *gin.Context) {
				clusters, _ := store.ListClusters()
				testResults, _ := store.ListTestResults("")

				c.JSON(200, gin.H{
					"clusters":     len(clusters),
					"test_results": len(testResults),
					"version":      "1.0.0",
				})
			})
		}

		// Provider simulation endpoints
		providerHandlers := NewProviderSimulationHandlers(store, logger)

		// Validation endpoints
		validate := v1.Group("/validate")
		{
			validate.GET("/:provider", providerHandlers.ValidateProvider)
		}

		// Provider simulation endpoints
		providers := v1.Group("/providers")
		{
			providers.POST("/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)
		}

		// Azure simulation endpoints (under /simulator)
		simulator := NewAzureHandlers(logger)

		sim := v1.Group("/simulator")
		{
			sim.POST("/azure/aks", simulator.SimulateAKS)
			sim.POST("/azure/budget", simulator.SimulateAzureBudget)
		}

		// Executor endpoints (under /executor)
		// exec := v1.Group("/executor")
		// These handlers should forward to the real cloud provider if simulation succeeded, or if --force is set
		// (Handlers to be implemented)
		// exec.POST("/azure/aks", executor.ExecuteAKS) // --dryrun, --force supported
		// exec.POST("/azure/budget", executor.ExecuteAzureBudget)
	}

	// Mock storage endpoints
	mocks := v1.Group("/mocks")
	{
		mocks.POST(":key", SaveMockData)
		mocks.GET(":key", GetMockData)
		mocks.POST(":key/save", SaveMockDataToFile)
		mocks.POST(":key/load", LoadMockDataFromFile)
	}

	// Documentation endpoint
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Documentation",
			"version": "v1",
			"endpoints": gin.H{
				"clusters": gin.H{
					"POST /api/v1/clusters": "Create a new AKS cluster",
				},
				"metrics": gin.H{
					"GET /api/v1/metrics/health": "Health check",
					"GET /api/v1/metrics/status": "Service status",
				},
				"providers": gin.H{
					"GET /api/v1/providers/:provider/info":                   "Get provider information",
					"GET /api/v1/providers/:provider/clusters":               "List clusters for provider",
					"POST /api/v1/providers/:provider/operations/:operation": "Simulate provider operation",
				},
				"validate": gin.H{
					"GET /api/v1/validate/:provider": "Validate provider configuration",
				},
				"simulator": gin.H{
					"POST /api/v1/simulator/azure/aks":    "Simulate AKS cluster creation",
					"POST /api/v1/simulator/azure/budget": "Simulate Azure budget",
				},
				"executor": gin.H{
					"POST /api/v1/executor/azure/aks":    "Execute AKS cluster creation (real cloud)",
					"POST /api/v1/executor/azure/budget": "Execute Azure budget (real cloud)",
				},
			},
		})
	})
}
