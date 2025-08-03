package api

import (
	"time"

	store "github.com/tronicum/punchbag-cube-testsuite/store"

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
			clusters.GET(":id", handlers.GetCluster)
			clusters.PUT(":id", handlers.UpdateCluster)
			clusters.DELETE(":id", handlers.DeleteCluster)

			// Test endpoints for specific clusters
			clusters.POST(":id/tests", handlers.RunTest)
			clusters.GET(":id/tests", handlers.ListTestResults)
		}

		// Test result endpoints
		tests := v1.Group("/tests")
		{
			tests.GET(":id", handlers.GetTestResult)
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


   providerSimHandlers := NewProviderSimulationHandlers(store, logger)

   // Simulation endpoints
   simulate := v1.Group("/simulate")
   {
	   simulate.POST("/providers/:provider/operations/:operation", providerSimHandlers.SimulateProviderOperation)
	   simulate.POST("/providers/:provider/buckets", providerSimHandlers.CreateSimulatedBucket)
	   simulate.GET("/providers/:provider/buckets", providerSimHandlers.ListSimulatedBuckets)
	   simulate.DELETE("/providers/:provider/buckets/:bucket", providerSimHandlers.DeleteSimulatedBucket)
	   // Generic AWS S3 simulation endpoint for SDK compatibility
	   simulate.Any("/aws-s3/*path", providerSimHandlers.GenericAWSS3SimHandler)
	   // Add more simulation endpoints as needed
   }

	   // Proxy endpoints (real provider, via server)
	   // proxy := v1.Group("/proxy")
	   // {
	   //     // TODO: Register proxy handlers here
	   //     // proxy.POST("/providers/:provider/operations/:operation", ...)
	   //     // proxy.DELETE("/providers/:provider/buckets/:bucket", ...)
	   // }

	   // direct := v1.Group("/direct")
	   // {
	   //     // TODO: Register direct handlers here
	   // }

	   // Validation endpoints (can be under simulate or proxy as appropriate)
		   validate := v1.Group("/validate")
		   {
				   validate.GET(":provider", providerSimHandlers.ValidateProvider)
		   }

	   // Azure simulation endpoints (legacy, to be migrated)
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
