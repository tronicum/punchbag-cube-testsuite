package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"yourapp/store" // replace with your actual import path
)

func NewMockLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestValidateProvider(t *testing.T) {
	// Set up Gin router
	router := gin.Default()
	logger := NewMockLogger() // Mock logger for testing
	providerHandlers := NewProviderSimulationHandlers(nil, logger)

	// Add validation endpoint
	router.GET("/validate/:provider", providerHandlers.ValidateProvider)

	// Test valid provider
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/validate/azure", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "azure")

	// Test invalid provider
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/validate/unknown", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "unsupported provider")
}

func TestSimulateProviderOperation(t *testing.T) {
	// Set up Gin router
	router := gin.Default()
	logger := NewMockLogger() // Mock logger for testing
	providerHandlers := NewProviderSimulationHandlers(nil, logger)

	// Add simulation endpoint
	router.POST("/providers/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)

	// Test valid operation
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/providers/azure/operations/create-cluster", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "create-cluster")

	// Test invalid operation
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/providers/unknown/operations/invalid", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "provider and operation are required")
}
