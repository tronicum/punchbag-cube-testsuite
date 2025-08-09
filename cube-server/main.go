package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	api "github.com/tronicum/punchbag-cube-testsuite/cube-server/api"
	"github.com/tronicum/punchbag-cube-testsuite/cube-server/internal"
	"github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
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
	// 1. Load server config (YAML or ENV)
	config, err := internal.LoadServerConfig()
	if err != nil {
		log.Printf("Could not load config, proceeding with defaults: %v", err)
	}
	// ENV/flag/config precedence for debug
	if config != nil && config.Debug {
		debugMode = true
	}
	if os.Getenv("CUBE_SERVER_DEBUG") == "1" {
		debugMode = true
	}
	// Initialize logger
	var logger *zap.Logger
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

	// Debug middleware: log all incoming requests (method, path, body)
	router.Use(func(c *gin.Context) {
		fmt.Printf("[GIN DEBUG] Incoming %s %s\n", c.Request.Method, c.Request.URL.Path)
		// Log all registered routes for debugging
		for _, ri := range router.Routes() {
			logger.Info("Registered route", zap.String("method", ri.Method), zap.String("path", ri.Path))
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			bodyBytes, _ := c.GetRawData()
			fmt.Printf("[GIN DEBUG] Request body: %s\n", string(bodyBytes))
			// Restore body for downstream handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		c.Next()
	})

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

	// 2. Create shared SimulationService with fastSimulate/debug from config/env
	fastSim := false
	if config != nil && config.FastSimulate {
		fastSim = true
	}
	if os.Getenv("FAST_SIMULATE") == "1" {
		fastSim = true
	}
	sim := simulation.NewSimulationServiceWithOptions(fastSim, debugMode)

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
	// Print incoming requests if debug enabled
	if debugMode {
		router.Use(func(c *gin.Context) {
			log.Printf("[SERVER DEBUG] %s %s", c.Request.Method, c.Request.URL.Path)
			c.Next()
		})
	}
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
