package main

import (
	   "log"
	   "net/http"
	   "os"

	   "github.com/gin-gonic/gin"
	   "go.uber.org/zap"

	   api "github.com/tronicum/punchbag-cube-testsuite/cube-server/api"
)

func main() {
	   // Check for --debug flag
	   debugMode := false
	   for _, arg := range os.Args {
			   if arg == "--debug" {
					   debugMode = true
					   break
			   }
	   }
	   // Initialize logger
	   var logger *zap.Logger
	   var err error
	   if debugMode {
			   logger, err = zap.NewDevelopment()
	   } else {
			   logger, err = zap.NewProduction()
	   }
	   if err != nil {
			   log.Fatalf("Failed to initialize logger: %v", err)
	   }
	   defer func() {
			   if err := logger.Sync(); err != nil {
					   log.Printf("Logger sync error: %v", err)
			   }
	   }()
	   if debugMode {
			   logger.Info("Debug mode enabled")
	   }

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


	// Register all API v1 routes (including simulation, proxy, etc.)
	api.SetupRoutes(router, nil, logger)

	// Start server
	// Get port from --port flag or environment variable
	port := "8080"
	if len(os.Args) > 1 {
		for i, arg := range os.Args {
			if arg == "--port" && i+1 < len(os.Args) {
				port = os.Args[i+1]
			}
		}
	}
	if envPort := os.Getenv("CUBE_SERVER_PORT"); envPort != "" {
		port = envPort
	}
	logger.Info("Starting Cube Server...", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}
