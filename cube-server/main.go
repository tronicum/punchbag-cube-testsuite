package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"punchbag-cube-testsuite/cube-server/api"
	"punchbag-cube-testsuite/cube-server/store"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize store
	dataStore := store.NewMemoryStore()

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
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "punchbag-cube-server",
		})
	})

	// Health check endpoint for automation (returns plain 'ok')
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Set up API routes
	api.SetupRoutes(router, dataStore, logger)

	// Improved validation endpoint: accepts provider/resource as query params, forwards to sim-server
	router.POST("/api/validation", func(c *gin.Context) {
		provider := c.Query("provider")
		resource := c.Query("resource")
		resp, err := http.PostForm("http://localhost:8080/api/validation?provider="+provider+"&resource="+resource, c.Request.PostForm)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Proxy endpoint: supports dryrun (validation) and real proxy
	addProxyRoute := func(path, target, provider, resource string) {
		u, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(u)
		router.Any(path, func(c *gin.Context) {
			if c.Query("dryrun") == "true" {
				// Forward to validation endpoint with provider/resource
				c.Request.URL.Path = "/api/validation"
				c.Request.Method = http.MethodPost
				c.Request.URL.RawQuery = "provider=" + provider + "&resource=" + resource
				router.HandleContext(c)
				return
			}
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Example: forward to simulation server (could be real cloud endpoint in production)
	addProxyRoute("/api/proxy/azure/aks", "http://localhost:8080/api/simulate/azure/aks", "azure", "aks")
	addProxyRoute("/api/proxy/azure/loganalytics", "http://localhost:8080/api/simulate/azure/loganalytics", "azure", "loganalytics")
	addProxyRoute("/api/proxy/azure/budget", "http://localhost:8080/api/simulate/azure/budget", "azure", "budget")
	addProxyRoute("/api/proxy/azure/appinsights", "http://localhost:8080/api/simulate/azure/appinsights", "azure", "appinsights")
	addProxyRoute("/api/proxy/aws/eks", "http://localhost:8080/api/simulate/aws/eks", "aws", "eks")
	addProxyRoute("/api/proxy/aws/s3", "http://localhost:8080/api/simulate/aws/s3", "aws", "s3")
	addProxyRoute("/api/proxy/gcp/gke", "http://localhost:8080/api/simulate/gcp/gke", "gcp", "gke")

	// Execute endpoint: forwards to sim-server for now, can be extended for real cloud
	router.POST("/api/execute/:provider/:resource", func(c *gin.Context) {
		provider := c.Param("provider")
		resource := c.Param("resource")
		// Forward to sim-server for now
		resp, err := http.PostForm("http://localhost:8080/api/simulate/"+provider+"/"+resource, c.Request.PostForm)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting cube-server", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
