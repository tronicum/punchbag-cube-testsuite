package main

import (
	   "log"
	   "net/http"
	   "os"

	   "github.com/gin-gonic/gin"
	   "go.uber.org/zap"

	   api "github.com/tronicum/punchbag-cube-testsuite/cube-server/api"
	   "github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
	   "github.com/tronicum/punchbag-cube-testsuite/cube-server/internal"
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



	// 1. Load server config (YAML or ENV)
	config, err := internal.LoadServerConfig()
	if err != nil {
		logger.Warn("Could not load config, proceeding with defaults", zap.Error(err))
	}

	// 2. Create shared SimulationService
	sim := simulation.NewSimulationService()

	// 3. Inject dummy buckets if needed (from config or ENV)
	if config != nil && config.Storage.DummyBuckets != nil {
		conv := map[string][]struct{ Name, Region string }{}
		for provider, buckets := range config.Storage.DummyBuckets {
			for _, b := range buckets {
				conv[provider] = append(conv[provider], struct{ Name, Region string }{b.Name, b.Region})
			}
		}
		internal.InjectDummyBucketsIfNeeded(sim, conv)
	} else {
		// fallback: try ENV/file-based config
		dummy := internal.LoadDummyBucketsConfig()
		if dummy != nil {
			conv := map[string][]struct{ Name, Region string }{}
			for provider, buckets := range dummy {
				for _, b := range buckets {
					conv[provider] = append(conv[provider], struct{ Name, Region string }{b.Name, b.Region})
				}
			}
			internal.InjectDummyBucketsIfNeeded(sim, conv)
		}
	}

	// 4. Register all API v1 routes (including simulation, proxy, etc.)
	api.SetupRoutes(router, nil, logger, sim)

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
