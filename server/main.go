package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tronicum/punchbag-cube-testsuite/server/api"
	"github.com/tronicum/punchbag-cube-testsuite/store"
	"github.com/tronicum/punchbag-cube-testsuite/server/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getLogger() (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	if config.IsDebug() {
		logCfg = zap.NewDevelopmentConfig()
		logCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logFile := os.Getenv("LOGFILE")
	if logFile == "" {
		logFile = "server.log"
	}
	logCfg.OutputPaths = []string{"stdout", logFile}
	return logCfg.Build()
}

func main() {
	// Initialize logger
	logger, err := getLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize store
	dataStore := store.NewMemoryStore()

	// Set up Gin router
	if config.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if config.IsDebug() {
		router.Use(gin.Logger())
	}
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
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "punchbag-cube-testsuite-server",
		})
	})

	// Set up API routes
	api.SetupRoutes(router, dataStore, logger)

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
