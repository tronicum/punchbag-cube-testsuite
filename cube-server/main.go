package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tronicum/punchbag-cube-testsuite/cube-server/sim"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Set up Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Simulation endpoints (migrated from sim-server)
	router.POST("/api/simulate/azure/aks", gin.WrapF(sim.HandleAks))
	router.GET("/api/simulate/azure/aks", gin.WrapF(sim.HandleAks))
	router.DELETE("/api/simulate/azure/aks", gin.WrapF(sim.HandleAks))

	router.POST("/api/simulate/azure/loganalytics", gin.WrapF(sim.HandleLogAnalytics))
	router.GET("/api/simulate/azure/loganalytics", gin.WrapF(sim.HandleLogAnalytics))

	router.POST("/api/validation", gin.WrapF(sim.HandleValidation))

	// Start server
	logger.Info("Starting Cube Server...")
	router.Run(":8080")
}
